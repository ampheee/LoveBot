FROM golang:latest

COPY . /go/src/app
WORKDIR /go/src/app/cmd
RUN go mod download

RUN apt-get update && apt-get install -y postgresql-client

EXPOSE 8080

RUN go build main.go

CMD ["./main"]