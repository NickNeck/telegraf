FROM artifactory.gcxi.de:8443/golang:latest
MAINTAINER gcx team <p+s@grandcentrix.net>

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -qy --no-install-recommends \
      python \
      ruby \
      ruby-dev \
      rubygems \
      build-essential \
      rpm-common \
      make

RUN gem install --no-ri --no-rdoc fpm

# Clean up
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go

CMD ["/bin/bash"]

