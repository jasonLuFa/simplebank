# build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.16
WORKDIR /app
# the dot mean /app in alphine:3.16
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
RUN chmod +x start.sh
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
# it will similar to running /app/start.sh /app/main command, but more flexibility to replace CMD in runtime
ENTRYPOINT [ "/app/start.sh" ]
CMD ["/app/main"]