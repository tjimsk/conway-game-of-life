FROM golang:1.11 as builder
COPY server/go.mod /root/tjimsk/game-of-life/go.mod
WORKDIR /root/tjimsk/game-of-life
RUN go mod download
COPY server /root/tjimsk/game-of-life
WORKDIR /root/tjimsk/game-of-life/main
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/game-of-life -a -installsuffix cgo .
FROM alpine:3.4 as image
COPY --from=builder /go/bin/game-of-life /go/bin/game-of-life
COPY client/dist /etc/game-of-life/static
ENV PORT=:80
ENV INTERVAL=1000
ENV HEIGHT=64
ENV WIDTH=120
ENV STATIC="/etc/game-of-life/static"
ENV SEED=false
EXPOSE 80
ENTRYPOINT /go/bin/game-of-life
