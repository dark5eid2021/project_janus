# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o riskassessor

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/riskassessor .
# by default, this app should run... For example:
# riskassessor /path/to/logfile /path/to/codechanges
CMD ["./riskassessor"]