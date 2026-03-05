FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/schema.sql .

EXPOSE 8000
CMD ["./main"]
