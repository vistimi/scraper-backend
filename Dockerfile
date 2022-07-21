FROM golang:1.18-alpine as builder

RUN apk add --no-cache git

WORKDIR /usr/tmp

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o scraper src/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /usr/tmp/scraper /usr/app/scraper

WORKDIR /usr/app

EXPOSE 8080

RUN echo $IMAGE_PATH
RUN echo $SCRAPER_DB

CMD ["./scraper"]