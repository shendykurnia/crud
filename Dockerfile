FROM        golang:1.7-alpine

# git
RUN apk add --no-cache git mercurial
# gcc
RUN apk add --update alpine-sdk

# migration
RUN mkdir /migration
ADD migration/ /migration
RUN go get -v github.com/mattes/migrate
CMD ["migrate", "-url postgres://tokopedia:tokopedia@postgres:5432/tokopedia", "-path /migration", "up"]

# config
RUN mkdir /appconfig
ADD config/ /appconfig

# govendor
RUN go get -v -u github.com/kardianos/govendor

# app
RUN mkdir /go/src/web
ADD src/web/ /go/src/web
WORKDIR /go/src/web
RUN govendor sync
RUN go build -o /bin/app app.go

CMD ["/bin/app", "-config=/appconfig"]

EXPOSE 9000
