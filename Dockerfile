# step 1 - loading dependences
FROM golang:1.20-alpine AS modules
WORKDIR /modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# step 2 - compiling the source code into binary
FROM golang:1.20-alpine AS build
COPY --from=modules /go/pkg /go/pkg
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY gen ./gen
COPY go.mod ./
COPY go.sum ./
RUN CGO_ENABLED=0 \
    go build -o /bin/app ./cmd/app

# step 3 - running binary application
FROM scratch
COPY --from=build /bin/app /app
COPY configs /configs
EXPOSE ${GRPC_PORT}
CMD ["/app"]