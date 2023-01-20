ARG VARIANT=golang:1.19.0-alpine
ARG RUNNER=workflow
ARG GO_ALPINE_VARIANT=golang:1.19.0-alpine

#-------------------------
#    GOLANG BUILDER
#-------------------------
FROM ${GO_ALPINE_VARIANT} as builder-alpine-go

#-------------------------
#    BUILDER FINAL
#-------------------------
FROM ${VARIANT} as builder-final

#-------------------------
#    RUNNER
#-------------------------
#                    --->   workflow   ---
#                   /                      \
#  builder-final ---                        ---> runner
#                   \                      /
#                    ---> devcontainer ---

#-------------------------
#    RUNNER WORKFLOW
#-------------------------
FROM builder-final AS builder-workflow

WORKDIR /usr/tmp

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
ENV GIN_MODE=release
RUN go build -o scraper src/main.go

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

COPY --from=builder-workflow /usr/tmp/scraper /usr/app/scraper

WORKDIR /usr/app

EXPOSE 8080

CMD ["./scraper"]

#-------------------------
#    RUNNER DEVCONTAINER
#-------------------------
FROM builder-final AS runner-devcontainer

#-------------------------
#       RUNNER
#-------------------------
FROM runner-${RUNNER} AS runner