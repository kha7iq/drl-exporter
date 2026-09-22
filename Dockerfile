FROM alpine:latest
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/drl-exporter"]
COPY $TARGETPLATFORM/drl-exporter /usr/bin/drl-exporter