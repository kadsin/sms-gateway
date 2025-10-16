FROM golang:1.25.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./build/balance-worker ./cmd/balance-worker

CMD ./build/balance-worker $CONSUMER_COUNT