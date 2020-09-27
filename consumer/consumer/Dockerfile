FROM golang:1.15-alpine AS build

RUN apk add --update --no-cache tzdata

WORKDIR /app
COPY . .
RUN go mod init consumer && CGO_ENABLED=0 go build -ldflags "-s -w"

FROM scratch

ENV TZ Europe/Moscow

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/blossom-observer /app/blossom-observer

# Use an unprivileged user.
USER 10001:10001

ENTRYPOINT ["/app/blossom-observer"]
#CMD ["sh", "-c", "tail -f /dev/null"]