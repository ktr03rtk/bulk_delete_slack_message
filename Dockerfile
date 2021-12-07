FROM golang:alpine

WORKDIR /go/src

COPY . .
RUN go mod download && go build -o ../bin ./...

ENTRYPOINT [ "/go/bin/bulkdelete" ]
