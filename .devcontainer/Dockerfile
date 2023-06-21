ARG ALPINE_VARIANT=mcr.microsoft.com/devcontainers/base:alpine-3.16
ARG GO_ALPINE_VARIANT=golang:1.19.0-alpine
ARG PYTHON_ALPINE_VARIANT=python:3.10.5-alpine

#-------------------------
#    BUILDER GOLANG
#-------------------------
FROM ${GO_ALPINE_VARIANT} as builder-alpine-go

#-------------------------
#    BUILDER PYTHON
#-------------------------
FROM $PYTHON_ALPINE_VARIANT as builder-alpine-python

ARG AWS_CLI_VERSION=2.9.0
RUN apk add --no-cache git unzip groff build-base libffi-dev cmake
RUN git clone --single-branch --depth 1 -b $AWS_CLI_VERSION https://github.com/aws/aws-cli.git

WORKDIR aws-cli
RUN python -m venv venv
RUN . venv/bin/activate
RUN scripts/installers/make-exe
RUN unzip -q dist/awscli-exe.zip
RUN aws/install --bin-dir /aws-cli-bin
RUN /aws-cli-bin/aws --version

# reduce image size: remove autocomplete and examples
RUN rm -rf \
    /usr/local/aws-cli/v2/current/dist/aws_completer \
    /usr/local/aws-cli/v2/current/dist/awscli/data/ac.index \
    /usr/local/aws-cli/v2/current/dist/awscli/examples
RUN find /usr/local/aws-cli/v2/current/dist/awscli/data -name completions-1*.json -delete
RUN find /usr/local/aws-cli/v2/current/dist/awscli/botocore/data -name examples-1.json -delete

#---------------
#    RUNNER
#---------------
FROM $ALPINE_VARIANT AS runner

# # non-root user
# # FIXME: use vscode user
# RUN apk add --no-cache shadow sudo
# ARG USERNAME=user
# ARG USER_UID=1001
# ARG USER_GID=$USER_UID
# RUN groupadd --gid $USER_GID $USERNAME \
#     && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
#     # Add sudo support. Omit if you don't need to install software after connecting.
#     && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
#     && chmod 0440 /etc/sudoers.d/$USERNAME
# # Automatically update the UID/GID
# RUN groupmod --gid $USER_GID $USERNAME \
#     && usermod --uid $USER_UID --gid $USER_GID $USERNAME \
#     && chown -R $USER_UID:$USER_GID /home/$USERNAME
# USER $USERNAME

# Golang
COPY --chown=$USERNAME:$USER_GID --from=builder-alpine-go /usr/local/go/ /usr/local/go/
COPY --chown=$USERNAME:$USER_GID --from=builder-alpine-go /go/ /go/
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH
RUN sudo mkdir -p "$GOPATH/src" "$GOPATH/bin" && sudo chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
RUN go version

# setup
RUN go install github.com/cweill/gotests/gotests@latest \
    && go install github.com/fatih/gomodifytags@latest \
    && go install github.com/josharian/impl@latest \
    && go install github.com/haya14busa/goplay/cmd/goplay@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/davidrjenni/reftools/cmd/fillstruct@latest

# aws cli
COPY --chown=$USERNAME:$USER_GID --from=builder-alpine-python /usr/local/aws-cli/ /usr/local/aws-cli/
COPY --chown=$USERNAME:$USER_GID --from=builder-alpine-python /aws-cli-bin/ /usr/local/bin/
RUN aws --version