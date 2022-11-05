# build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
# the dot mean /app in alphine:3.13
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]