FROM golang:1.16.5

WORKDIR /avitoChatAPI

COPY . .

RUN go mod download
