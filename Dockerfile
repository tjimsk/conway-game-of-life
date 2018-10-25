FROM golang:latest

RUN go get "github.com/gorilla/websocket"
RUN go get "github.com/spf13/viper"

WORKDIR /go/src/client/dist
COPY client/dist .

WORKDIR /go/src/server
COPY server .

RUN go build -o conway *.go

ENV PORT=:8080
ENV HEIGHT=100
ENV WIDTH=150

EXPOSE 8080

CMD ["./conway"]
