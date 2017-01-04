FROM golang:1.7

COPY . /go/src/github.com/billybobjoeaglt/website/
WORKDIR /go/src/github.com/billybobjoeaglt/website/
RUN go build

EXPOSE 80 80
#docker run -d -p 80:80 --name website website -port 80 -prod
ENTRYPOINT ["/go/src/github.com/billybobjoeaglt/website/website"]
CMD ["-port 80 -prod"]