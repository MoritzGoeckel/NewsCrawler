# Downloader
#
# VERSION               0.1

FROM golang:1.8

WORKDIR /go/src/app

# COPY . .
COPY ./content/ /go/src/app/

# RUN go get github.com/PuerkitoBio/goquery
# RUN go get -u github.com/go-redis/redis

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]

LABEL name="downloader"

# RUN find /usr/local/apache2/htdocs/ -exec chmod 775 {} \;
