# create a new image with golang
FROM golang:1.18-alpine as builder
RUN apk add --no-cache git

WORKDIR /usr/tmp

# install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# build project
COPY . .
RUN go build -o scraper src/main.go

# create a new empty image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# copy the build file
COPY --from=builder /usr/tmp/scraper /usr/app/scraper

WORKDIR /usr/app

EXPOSE 8080

CMD ["./scraper"]