FROM golang:1.25.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./build/wallet-worker ./cmd/wallet-worker

CMD ./build/wallet-worker $CONSUMER_COUNT