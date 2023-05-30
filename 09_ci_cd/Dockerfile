FROM golang:1.18-alpine

COPY . /go/src/note_app

WORKDIR /go/src/note_app
RUN apk add make && make build
ENTRYPOINT [ "/go/src/note_app/bin/app" ]

EXPOSE 80/tcp
CMD [ "app" ]