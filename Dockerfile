FROM golang:1.17

RUN apt-get update && apt-get install -y \
  ffmpeg \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build bot.go

ENV TOKEN_TELEGRAM=$TOKEN_TELEGRAM

CMD ["./bot"]