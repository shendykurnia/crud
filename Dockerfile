FROM        golang:1.7-alpine

RUN mkdir /app

ADD src/vendor/ /usr/local/go/src/
ADD src/ /app

WORKDIR /app
RUN go build -o /bin/app app.go

CMD ["/bin/app"]

EXPOSE 9000
