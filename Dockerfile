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
ENV INTERVAL=600
ENV HEIGHT=75
ENV WIDTH=90
ENV STATIC="/etc/game-of-life/static/production"
ENV SEED=false
EXPOSE 80
ENTRYPOINT /go/bin/game-of-life
