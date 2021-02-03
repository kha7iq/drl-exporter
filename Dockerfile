# Build

FROM golang:1.15-alpine as build
ENV CGO_ENABLED=0 GOOS=linux
RUN mkdir /app
COPY . /app
WORKDIR /app/cmd/drl-exporter
RUN apk update && apk add --no-cache gcc musl-dev git

RUN go mod download

RUN  go build -ldflags '-w -s' -a  -o /app/bin/drl-exporter

# Final image
FROM alpine:latest

COPY --from=build /app/bin/drl-exporter /drl-exporter

CMD ["/drl-exporter"]