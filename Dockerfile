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

# RUN apk update
# RUN rc-service docker start
# # RUN sudo usermod -aG docker $USER && newgrp docker
# RUN docker ps

#-------------------------
#    RUNNER
#-------------------------
#                    ---> runner-${ALPINE_VARIANT}           ---
#                   /                                            \
#  builder-final ---                                              ---> runner
#                   \                                            /
#                    ---> runner-${DEVCONTAINTER_VARIANT}    ---

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

ARG USERNAME=user
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apk update && apk add --update sudo

RUN addgroup --gid $USER_GID $USERNAME \
    && adduser --uid $USER_UID -D -G $USERNAME $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME

COPY --from=builder-workflow /usr/tmp/scraper /usr/app/scraper

WORKDIR /usr/app

EXPOSE 8080
EXPOSE 27017

CMD ["./scraper"]

#-------------------------
#    RUNNER DEVCONTAINER
#-------------------------
FROM builder-final AS runner-devcontainer

#-------------------------
#       RUNNER
#-------------------------
FROM runner-${RUNNER} AS runner