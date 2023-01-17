# Build

FROM golang:1.19-alpine as build
ENV CGO_ENABLED=0 GOOS=linux
RUN apk update && apk add --no-cache gcc musl-dev git
RUN mkdir /app
COPY . /app
WORKDIR /app/cmd/drl-exporter

# RUN go mod download

RUN  go build -mod vendor -ldflags '-w -s' -a  -o /app/bin/drl-exporter

# Final image
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/bin/drl-exporter /drl-exporter

CMD ["/drl-exporter"]