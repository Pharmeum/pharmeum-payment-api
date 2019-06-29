FROM golang:1.12.1-stretch

WORKDIR $GOPATH/src/github.com/Pharmeum/pharmeum-payment-api/

COPY . .

RUN go build -o payment-api -v ./cmd/main.go