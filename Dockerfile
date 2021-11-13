# syntax=docker/dockerfile:1.2

FROM golang:1.16-alpine as puredns
ENV GO111MODULE on
RUN apk add --no-cache \
      git && \
    git clone https://github.com/d3mondev/puredns.git /opt/puredns

WORKDIR /opt/puredns
RUN go install ./...

FROM alpine:3.14 as massdns
RUN apk add --update --no-cache  \
      build-base \
      git \
      ldns \
      ldns-dev && \
    git clone --branch=master --depth=1 \
      https://github.com/blechschmidt/massdns.git && \
    cd massdns && \
    make

FROM alpine:3.14 as final
COPY --from=massdns /massdns/bin/massdns /usr/bin/massdns
COPY --from=puredns /go/bin/puredns /usr/bin/puredns

ENTRYPOINT [ "/usr/bin/puredns" ]