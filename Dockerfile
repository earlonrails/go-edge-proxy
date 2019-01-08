FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go get -u github.com/golang/dep/...
RUN dep ensure
RUN go build .

CMD ["./app"]
