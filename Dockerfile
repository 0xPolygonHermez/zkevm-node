# CONTAINER FOR BUILDING BINARY
FROM golang:1.16 AS build

ENV CGO_ENABLED=1

# INSTALL DEPENDENCIES
RUN go get -u github.com/gobuffalo/packr/v2/packr2
COPY go.mod go.sum /src/
RUN cd /src && go mod download

# BUILD BINARY
COPY . /src
RUN cd /src/db && packr2
RUN cd /src && make build

# CONTAINER FOR RUNNING BINARY
FROM golang:1.16
WORKDIR /app
COPY --from=build /src/dist/hezcore /app/hezcore
EXPOSE 8123
CMD ["./hezcore", "run"]