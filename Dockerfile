FROM alpine:3.20.2
ENTRYPOINT ["/usr/bin/drl-exporter"]
COPY drl-exporter /usr/bin/drl-exporter