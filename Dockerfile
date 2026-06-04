# Build stage
FROM golang:1.26 AS builder

WORKDIR /go/src/github.com/UKHomeOffice/acp-opensearch-alias-exporter

COPY go.mod ./

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go install -v \
            github.com/UKHomeOffice/acp-opensearch-alias-exporter

# Runtime stage
FROM alpine:3.22.4

RUN apk --no-cache upgrade
RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -S app -G app

USER 1000

COPY --from=builder --chown=app:app /go/bin/acp-opensearch-alias-exporter /acp-opensearch-alias-exporter
CMD ["/acp-opensearch-alias-exporter"]

