FROM golang:1.25-alpine AS build

RUN apk update && apk upgrade

WORKDIR /meshtastic2mastodon

COPY *.go go.sum go.mod /meshtastic2mastodon/

COPY protobufs /meshtastic2mastodon/protobufs

RUN go build -o /meshtastic2mastodon/meshtastic2mastodon

FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /meshtastic2mastodon

COPY --from=build /meshtastic2mastodon/meshtastic2mastodon /meshtastic2mastodon/meshtastic2mastodon

CMD ["/meshtastic2mastodon/meshtastic2mastodon"]
