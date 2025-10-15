FROM golang:1.25.2-alpine

WORKDIR /app

# ENV GOPROXY "https://goproxy.cn"
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./build/database-migrator-cli ./cmd/database-migrator-cli
RUN go build -o ./build/analytics-migrator-cli ./cmd/analytics-migrator-cli

CMD ./build/database-migrator-cli up && ./build/database-migrator-cli up