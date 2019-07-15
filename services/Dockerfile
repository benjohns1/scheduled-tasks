FROM golang:1.12.2 AS builder
COPY . /app/
WORKDIR /app/cmd/srv
RUN env GOOS=linux GOARCH=386 go build

FROM scratch
COPY --from=builder /app/cmd/srv/srv /
ENTRYPOINT ["/srv"]