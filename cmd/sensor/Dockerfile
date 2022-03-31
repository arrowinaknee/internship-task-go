FROM golang:1.17 as builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./...

FROM alpine:latest
COPY --from=builder /usr/local/bin/app ./
CMD ["./app"]