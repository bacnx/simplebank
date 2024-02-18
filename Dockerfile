# Build state
FROM  golang:1.22.0-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main ./main

CMD ["/app/main"]
EXPOSE 8080
