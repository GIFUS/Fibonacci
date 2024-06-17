FROM golang:latest

WORKDIR /app

COPY server/go.mod server/go.sum ./

RUN go mod download

COPY server/*.go ./

RUN go build -o server .

EXPOSE 8080

CMD ["./server"]
