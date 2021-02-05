FROM golang:1.15.7-alpine

ENV CGO_ENABLED=0

WORKDIR /app
COPY . .

CMD go test -v ./...
