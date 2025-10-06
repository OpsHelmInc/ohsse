FROM golang AS build

ENV CGO_ENABLED=0
RUN apt-get update
RUN apt-get upgrade -y
RUN update-ca-certificates --verbose

WORKDIR /build/cmd

COPY . /build

RUN go build -o ohsse

FROM scratch AS final

WORKDIR /ohsse
COPY --from=build /build/cmd/ohsse .
COPY --from=build /etc/ssl/ /etc/ssl/

ENTRYPOINT ["./ohsse"]