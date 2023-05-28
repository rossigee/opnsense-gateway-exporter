FROM golang:1.20 as builder
ADD . /build
WORKDIR /build
RUN go vet ./...
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o /build/opnsense-metrics

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/opnsense-metrics /bin/opnsense-metrics
EXPOSE 9576
ENTRYPOINT ["/bin/opnsense_exporter"]
