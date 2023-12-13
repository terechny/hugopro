
FROM golang:1.16-alpine as builder

WORKDIR /app

COPY . .

RUN go mod init chimiddlware && go mod tidy
RUN go build -o main

FROM alpine:latest

COPY --from=builder /app/main /main

EXPOSE 8080

CMD ["/main"]
