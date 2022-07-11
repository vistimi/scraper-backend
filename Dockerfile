FROM golang:1.18-alpine as builder

RUN apk add --no-cache git

WORKDIR /usr/src/tmp

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o scraper src/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /usr/src/tmp/scraper /usr/src/app/scraper

WORKDIR /usr/src/app

EXPOSE 8080

CMD ["./scraper"]