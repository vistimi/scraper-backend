ARG VARIANT=golang:1.19.0-alpine
ARG RUNNER=workflow

FROM ${VARIANT} as builder-final

# builder
FROM builder-final AS builder-workflow

WORKDIR /usr/tmp

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
ENV GIN_MODE=release
RUN go build -o scraper src/main.go

# runner
FROM builder-final AS runner-workflow

ARG USER_NAME=user
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN apk update && apk add --update sudo
RUN addgroup --gid $USER_GID $USER_NAME \
    && adduser --uid $USER_UID -D -G $USER_NAME $USER_NAME \
    && echo $USER_NAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USER_NAME \
    && chmod 0440 /etc/sudoers.d/$USER_NAME
USER $USER_NAME

WORKDIR /usr/app

COPY --from=builder-workflow /usr/tmp/scraper /usr/app/scraper
COPY --from=builder-workflow /usr/tmp/config /usr/app/config

EXPOSE 8080

CMD ["./scraper"]