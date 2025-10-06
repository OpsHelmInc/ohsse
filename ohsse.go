package ohsse

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"golang.org/x/net/http2"
)

// GetSSEStream configures and connects to a stream
func GetSSEStream(url string, authKey string) (io.ReadCloser, error) {
	// Default TLS config
	tls := &tls.Config{}
	// KEYLOG sets a file to save a keylog to
	// (https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS/Key_Log_Format).
	// This saves the current cryptographic keys used by the TLS connection to a
	// file so that other applications, for example wireshark:
	// https://wiki.wireshark.org/TLS#key-log-format can decrypt the TLS packets
	// in real time and provide useful introspection.  If the KEYLOG env var is
	// not set this has no effect
	if keyLogFile, ok := os.LookupEnv("KEYLOG"); ok {
		f, err := os.OpenFile(keyLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Printf("could not open keylog file=[%s], err=[%s].  Continuing without key log", keyLogFile, err)
		} else {
			defer func() {
				err = f.Close()
				if err != nil {
					log.Printf("could not close keylog file, err=[%s]", err)
				}
			}()
			tls.KeyLogWriter = f
			log.Printf("using file=[%s] as keylog file", keyLogFile)
		}
	}

	// Use HTTP2 transport, and can override defaults here.  If we set timeouts
	// then SSE will not work properly and connections will "fail", when they
	// have not failed
	tr := &http2.Transport{
		TLSClientConfig: tls,
		IdleConnTimeout: 0,
	}

	// Use the native http client to form a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Auth here if APIKEY env var is set
	if len(authKey) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", authKey))
	}

	// "*" is probably also acceptable, but protecting against issues if other types may be shared in the future
	req.Header.Set("Accept", "text/event-stream")
	// This needs to be here or else GO will default to gzip, which doesn't work
	// well with streaming and may break things. default uses zlib
	req.Header.Set("accept-encoding", "deflate")
	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// WatchStream watches a connected SSE stream, parses the incoming data and passes it to the supplied handler
func WatchStream(stream io.ReadCloser, streamHandler StreamHandler, eventHandler CloudEventHandler) error {
	// Eventbridge has a 256kb max event size, so we want to be using larger
	// than that, plus some wiggle room, for multiple events
	const bufSize = 256 * 1024 * 5
	var (
		entry SSE_Entry
	)

	// Use a scanner to read the body
	scanner := bufio.NewScanner(stream)
	buf := make([]byte, 0, bufSize)
	scanner.Buffer(buf, bufSize)

	// Main loop that keeps checking for new data being sent over the stream
	for scanner.Scan() {
		event := string(scanner.Bytes())
		// Handle each line
		for _, field := range strings.Split(event, "\n") {
			if len(field) == 0 {
				// An empty line (i.e. a '\n\n') is the sign of the end of and entry
				// Pass to handler now we have a complete entry
				go streamHandler(entry)

				// If this entry is of type event, unmarshal it into a cloudevent and pass to that handler
				if entry.Event == "data" {
					cEvent := cloudevents.NewEvent()
					err := json.Unmarshal([]byte(entry.Data), &cEvent)
					if err == nil {
						go eventHandler(cEvent)
					}
				}

				// Clear the entry out
				entry = SSE_Entry{}
				continue
			}
			dataType := strings.SplitN(field, ":", 2)[0]
			if len(field) <= len(dataType)+1 {
				continue
			}
			dataPayload := strings.TrimSpace(field[len(dataType)+1:])
			switch dataType {
			case "":
				entry.Comment = field
			case "data":
				entry.Data = dataPayload
			case "event":
				entry.Event = dataPayload
			case "retry":
				entry.Retry = dataPayload
			case "id":
				entry.ID = dataPayload
			default:
				// Rather than lose unknown fields, store them up in this map
				if entry.Unknown == nil {
					entry.Unknown = make(map[string]string)
				}
				entry.Unknown[dataType] = dataPayload
			}
		}
	}
	return scanner.Err()
}
