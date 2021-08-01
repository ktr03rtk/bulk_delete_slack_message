FROM golang:alpine

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./

RUN go build -o ../bin ./...

ENTRYPOINT [ "/go/bin/bulkdelete" ]
