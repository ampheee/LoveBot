FROM golang:alpine

WORKDIR ./LoveBot
COPY . .
RUN go mod download
EXPOSE 8080
RUN go build main.go

ENTRYPOINT ["./main"]