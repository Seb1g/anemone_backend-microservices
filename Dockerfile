FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG SERVICE_NAME

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./internal/${SERVICE_NAME}/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /main .

CMD ["./main"]