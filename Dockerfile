# create a new image with golang
FROM golang:1.19.2-alpine as builder

WORKDIR /usr/tmp

# install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# build project
COPY . .
ENV GIN_MODE=release
RUN go build -o scraper src/main.go

# create a new empty image
FROM alpine:latest

# copy the build file
COPY --from=builder /usr/tmp/scraper /usr/app/scraper

WORKDIR /usr/app

# port for scraper
EXPOSE 8080

# port for mongodb
EXPOSE 27017

CMD ["./scraper"]