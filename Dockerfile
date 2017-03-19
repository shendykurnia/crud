FROM        golang:1.7-alpine

# git
RUN apk add --no-cache git mercurial
# gcc
RUN apk add --update alpine-sdk

# migration
# not needed anymore since i decided not to use postgres for this project completion
#RUN mkdir /migration
#ADD migration/ /migration
#RUN go get -v github.com/mattes/migrate
#CMD ["migrate", "-url postgres://tokopedia:tokopedia@postgres:5432/tokopedia", "-path /migration", "up"]

# config
# not needed anymore since i decided not to use postgres, redis, and nsq for this project completion
#RUN mkdir /appconfig
#ADD config/ /appconfig

# govendor
RUN go get -v -u github.com/kardianos/govendor

# app
COPY src/ /go/src
WORKDIR /go/src/web
RUN govendor sync

WORKDIR /go/src/myrouter
RUN govendor sync
RUN go test

WORKDIR /go/src/mymodel
RUN govendor sync
RUN go test

WORKDIR /go/src/web
RUN go test
RUN go build -o /bin/app app.go

#CMD ["/bin/app", "-config=/appconfig"]
CMD ["/bin/app"]

EXPOSE 9000
