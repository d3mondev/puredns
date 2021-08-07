# Specify the official Ubuntu latest base image
FROM ubuntu

# Set a working directory
WORKDIR /usr/src

# for apt to be noninteractive
ENV DEBIAN_FRONTEND noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN true

## preesed tzdata, update package index, upgrade packages and install needed software
RUN truncate -s0 /tmp/preseed.cfg; \
    echo "tzdata tzdata/Areas select Europe" >> /tmp/preseed.cfg; \
    echo "tzdata tzdata/Zones/Europe select Berlin" >> /tmp/preseed.cfg; \
    debconf-set-selections /tmp/preseed.cfg && \
    rm -f /etc/timezone /etc/localtime && \
    apt-get update && \
    apt-get install -y tzdata

## cleanup of files from setup
RUN rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install dependencies
RUN apt update && \
    apt install -y make git libsqlite3-dev libxslt-dev libxml2-dev zlib1g-dev gcc wget && \
    apt clean

# Install MassDNS
RUN git clone https://github.com/blechschmidt/massdns.git /usr/src/massdns && \
    cd /usr/src/massdns && \
    make && \
	cp /usr/src/massdns/bin/massdns /usr/local/bin/massdns

# Install Golang
ENV GOROOT=/usr/local/go 
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
RUN cd /tmp && \
	wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz && \
	tar -xvf go1.16.4.linux-amd64.tar.gz && \
	mv go /usr/local && \
	GO111MODULE=on go get github.com/d3mondev/puredns/v2

# Default command
CMD ["/bin/bash"]