FROM golang:latest

WORKDIR /app

RUN apt-get update && apt-get install -y libssl-dev pkg-config

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/main.go

EXPOSE 5005

CMD ["./server"]
