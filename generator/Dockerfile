FROM golang:latest

WORKDIR /app

COPY generator/go.mod generator/go.sum ./

RUN go mod download

COPY generator/*.go ./

RUN go build -o generator .

EXPOSE 8080

CMD ["./generator"]
