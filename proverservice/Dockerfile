FROM golang:1.17

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build -o main .

EXPOSE 50051

CMD ["/app/main"]