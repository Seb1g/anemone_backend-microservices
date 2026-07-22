# --- STAGE 1: Builder ---
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/main ./cmd/main.go

# --- STAGE 2: Final Run ---
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /home/appuser

COPY --from=builder /app/main .

USER appuser

EXPOSE 8080

CMD ["./main"]