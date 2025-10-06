package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/OpsHelmInc/localdevtools/ohsse"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	var (
		url string
		key string
	)
	flag.StringVar(&key, "key", "", "api key to use for connection")
	flag.StringVar(&url, "endpoint", "https://streaming.opshelm.com/v1/consume/stream", "endpoint to use")
	flag.Parse()

	if len(key) == 0 {
		key = os.Getenv("APIKEY")
	}

	for {
		stream, err := ohsse.GetSSEStream(url, key)
		if err != nil {
			fmt.Printf("could not set up SSE stream: %s\n", err)
			os.Exit(-1)
		}
		err = ohsse.WatchStream(stream, myStreamHandler, myCloudEventHandler)
		if err != nil {
			// We don't need to check the error on this, as an error in closing
			// the stream that's errored is moot for the purposes of this and
			// we're going to restart anyway
			_ = stream.Close()
			// Sleep before resuming to make sure that everything has closed
			// down properly first
			time.Sleep(10 * time.Second)
		}
	}
}

// This is just a generic handler that handles events of all types.  Most likely
// an end user isn't likely to want to use this as most event types (other than
// data) pertain to the internals of SSE, not the message we're delivering
func myStreamHandler(entry ohsse.SSE_Entry) {
	// Not doing anything special here.  But we can use it to handle PINGs and
	// data separately, etc.
}

// This is a handler that receives events of type data that have been
// unmarhshalled into a cloudevent.
func myCloudEventHandler(ce cloudevents.Event) {
	// Simply outputting events one per line
	fmt.Printf("%s\n", ce.Data())
}
