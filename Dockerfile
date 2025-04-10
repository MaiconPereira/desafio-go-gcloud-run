FROM golang:1.22-alpine AS builder

RUN apk --no-cache add ca-certificates && \
    addgroup -S app && \
    adduser -S -G app app

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

FROM scratch

WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /server /server

USER app
EXPOSE 8080
ENTRYPOINT ["/server"]
