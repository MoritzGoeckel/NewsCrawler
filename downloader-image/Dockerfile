# Downloader
#
# VERSION               0.1

FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]

LABEL name="downloader"

COPY ./content/ /go/src/app/

# RUN find /usr/local/apache2/htdocs/ -exec chmod 775 {} \;
