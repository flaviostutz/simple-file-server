FROM golang:1.14.0-alpine3.11

ENV LOG_LEVEL 'info'
ENV WRITE_SHARED_KEY ''
ENV READ_SHARED_KEY ''
ENV LOCATION_BASE_URL ''
ENV DATA_DIR '/data'

RUN apk add --no-cache build-base

WORKDIR /simple-file-server

ADD go.mod .
RUN go mod download

ADD / /simple-file-server
WORKDIR /simple-file-server
# RUN go test
RUN go build -o /usr/bin/simple-file-server

ADD startup.sh /

VOLUME [ "/data" ]

CMD [ "/startup.sh" ]

