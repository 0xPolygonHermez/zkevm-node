FROM golang:1.17

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .

EXPOSE 50051

CMD ["/app/main"]