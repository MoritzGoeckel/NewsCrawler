FROM golang:1.8

WORKDIR /go/src/app

# COPY . .
COPY ./content/ /go/src/app/

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]

LABEL name="api-server"
