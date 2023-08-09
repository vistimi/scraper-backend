ARG GO_ALPINE_VARIANT=golang:1.19.0-alpine
ARG VARIANT=alpine:3.16
ARG PORT=3000

# builder
FROM $VARIANT AS builder

RUN apk add --update --no-cache go

WORKDIR /usr/tmp

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./config ./config

ARG TARGETOS TARGETARCH
ENV GOPATH /go
ENV USERNAME root
ENV GIN_MODE=release
RUN --mount=target=. \
    --mount=type=cache,target=/$USERNAME/.cache/go-build \
    --mount=type=cache,target=$GOPATH/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o /out/scraper src/main.go

# runner
FROM $VARIANT AS runner

RUN apk add --update --no-cache shadow go

RUN apk add 
ARG USERNAME=user
ARG USER_UID=1001
ARG USER_GID=$USER_UID
RUN addgroup --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME
# # Add sudo support. Omit if you don't need to install software after connecting.
# && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
# && chmod 0440 /etc/sudoers.d/$USERNAME
USER $USERNAME

WORKDIR /usr/app
COPY --chown=$USERNAME:$USER_GID --from=builder /out/scraper ./
COPY --chown=$USERNAME:$USER_GID --from=builder /usr/tmp/config/config.yml ./config/config.yml

EXPOSE $PORT

CMD ["./scraper"]