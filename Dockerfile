FROM alpine:latest
ENTRYPOINT ["/usr/bin/drl-exporter"]
COPY drl-exporter /usr/bin/drl-exporter