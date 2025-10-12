FROM golang:1.25.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./build/server ./cmd/server

RUN chmod +x ./deployments/docker-entrypoint.sh
RUN mv ./deployments/docker-entrypoint.sh /usr/bin

RUN apk add --no-cache make

ENTRYPOINT docker-entrypoint.sh

EXPOSE 3000
