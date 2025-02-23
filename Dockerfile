FROM golang:1.23.3

WORKDIR /app

COPY go.* /app/
RUN go mod download

COPY . . 

RUN cd ./cmd/api && go build -o main .

CMD ["./cmd/api/main"]