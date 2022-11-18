
FROM amazon/aws-cli:latest as builder

# WORKDIR /github/workflow/

# RUN  curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
#     unzip -q awscliv2.zip && sudo ./aws/install

# RUN sudo apt-get update && sudo apt-get -q install -yy less

FROM alpine:latest

COPY --from=builder /usr/local/bin/aws /usr/local/bin/aws
COPY --from=builder /usr/bin/less /usr/bin/less