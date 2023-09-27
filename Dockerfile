FROM golang:1.20 as builder

WORKDIR /go/src/gitlab.digital.homeoffice.gov.uk/acp/opensearch-reporter

COPY go.mod ./

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go install -v \
            gitlab.digital.homeoffice.gov.uk/acp/opensearch-reporter

FROM alpine:3.17
RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -S app -G app

USER 1000

COPY --from=builder /go/bin/opensearch-reporter /opensearch-reporter
CMD ["/opensearch-reporter"]

