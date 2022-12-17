ARG VARIANT=alpine:3.16
ARG RUNNER=workflow
ARG ALPINE_VARIANT=alpine:3.16
ARG GO_ALPINE_VARIANT=golang:1.19.0-alpine

#-------------------------
#    GOLANG BUILDER
#-------------------------
FROM ${GO_ALPINE_VARIANT} as builder-alpine-go

#-------------------------
#    BUILDER FINAL
#-------------------------
FROM ${VARIANT} as builder-final

RUN apk update

# RUN rc-service docker start
# # RUN sudo usermod -aG docker $USER && newgrp docker
# RUN docker ps

# Golang
COPY --from=builder-alpine-go /usr/local/go/ /usr/local/go/
COPY --from=builder-alpine-go /go/ /go/
# ENV GOROOT /go
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
RUN go version

RUN go install github.com/cweill/gotests/gotests@latest \
    && go install github.com/fatih/gomodifytags@latest \
    && go install github.com/josharian/impl@latest \
    && go install github.com/haya14busa/goplay/cmd/goplay@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install golang.org/x/tools/gopls@latest

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
ENV RUNNER=$RUNNER