# build stage
FROM golang:alpine as build-stage

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./

RUN go build -o ../bin ./...

# ENTRYPOINT [ "/bin/sh" ]
ENTRYPOINT [ "go", "run", "cmd/bulkdelete/main.go" ]

# runtime stage
# FROM alpine as runtime-stage

# COPY --from=build-stage /go/bin /go/bin

# ENTRYPOINT [ "/go/bin/bulkdelete" ]
