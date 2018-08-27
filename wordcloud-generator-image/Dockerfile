FROM golang:1.10

WORKDIR /go/src/app
COPY ./content/ .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]

LABEL name="downloader"
