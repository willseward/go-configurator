# Docker Auditd Exporter - Host Agent
# Wills Ward
# Texas Christian University
# 03/31/16

FROM golang:1.6

MAINTAINER William Ward <wills.e.ward@tcu.edu>

ONBUILD ADD ./templates /var/go-configurator/templates

# Build binary
ADD . /go/src

RUN go get -v go-configurator
RUN go install -v go-configurator

EXPOSE 80

ENTRYPOINT ["/go/bin/go-configurator"]
