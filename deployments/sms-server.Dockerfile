FROM golang:1.25.2-alpine

WORKDIR /app

# ENV GOPROXY "https://goproxy.cn"
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./build/server ./cmd/server

ENTRYPOINT ./build/server

EXPOSE 3000
