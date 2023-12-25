FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache bash

RUN go mod download


WORKDIR /app/cmd/
RUN go build -o app

CMD ["./app"]