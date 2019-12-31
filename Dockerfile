FROM golang:1.13-alpine

WORKDIR /go/src/app
COPY gofamily.go .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]