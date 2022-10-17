# CONTAINER FOR BUILDING BINARY
FROM golang:1.17 AS build

ENV CGO_ENABLED=0

# INSTALL DEPENDENCIES
RUN go install github.com/gobuffalo/packr/v2/packr2@v2.8.3
COPY go.mod go.sum /src/
RUN cd /src && go mod download

# BUILD BINARY
COPY . /src
RUN cd /src/db && packr2
RUN cd /src && make build

# CONTAINER FOR RUNNING BINARY
FROM alpine:3.16.0
COPY --from=build /src/dist/zkevm-node /app/zkevm-node
COPY --from=build /src/config/config.local.toml /app/config.local.toml
EXPOSE 8123
CMD ["/bin/sh", "-c", "/app/zkevm-node run"]