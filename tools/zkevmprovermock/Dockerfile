FROM golang:1.17-alpine AS build

ENV CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base && \
    go build -ldflags '-extldflags "-static"' -o ./zkprover-mock ./tools/zkevmprovermock/cmd/...

FROM alpine:3.16.0
COPY --from=build /app/zkprover-mock /app/zkprover-mock
CMD ["/bin/sh", "-c", "/app/zkprover-mock server"]
