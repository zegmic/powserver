FROM golang:1.21 AS builder
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  appuser
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /powserver cmd/server/main.go

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /powserver /powserver
EXPOSE 8080
USER appuser:appuser
ENTRYPOINT ["/powserver"]