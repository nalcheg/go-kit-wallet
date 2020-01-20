FROM golang:1.13.6-alpine3.11 AS builder

RUN apk add --update ca-certificates

ADD . /build/

RUN cd /build && CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o /wallet ./cmd/main.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /wallet /app/wallet
ADD ./migrations /app/migrations

WORKDIR /app

ENTRYPOINT ["/app/wallet"]
