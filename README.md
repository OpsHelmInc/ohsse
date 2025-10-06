# ohsse

OpsHelm Streaming API client

This repository includes the `ohsse` general use package for building opshelm streaming API clients, and a simple client using this package in the `cmd` directory.  Note this client does not output any debugging data (with the exception of a failure during setup), it simply outputs one event per line in JSON format.

## Installation

There are two main options for installing and using `ohsse`

### Docker

This repository has a dockerfile included and so can you can:

- checkout this repo: `git clone https://github.com/OpsHelmInc/ohsse.git`
- change the `cmd` directory: `cd ohsse`
- build the docker image: `docker build .`
- run the docker image with your api key: `docker run <imageid> -key <apikey>`

### Using an existing golang environment

If you already have a golang environment setup, you can compile the source locally

- checkout this repo: `git clone https://github.com/OpsHelmInc/ohsse.git`
- change the `cmd` directory: `cd ohsse`
- build the binary: `go build . -out ohsse`
- run the binary locally with your api key: `./ohsse -key <apikey>`
