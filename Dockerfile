# Docker Auditd Exporter - Host Agent
# Wills Ward
# 03/31/16

FROM golang:1.6

MAINTAINER William Ward <wills.e.ward@tcu.edu>

ONBUILD ADD ./templates /var/go-configurator/templates

# Build binary
RUN go get github.com/willseward/go-configurator
RUN go get -v ./...
RUN go install -v github.com/willseward/go-configurator
EXPOSE 80

ENTRYPOINT ["/go/bin/go-configurator"]
