FROM golang:1.23.3-alpine3.20

ENV LOG_LEVEL='info' \
    WRITE_SHARED_KEY='' \
    READ_SHARED_KEY='' \
    LOCATION_BASE_URL='' \
    DATA_DIR='/data' \
    PORT=''

RUN apk add --no-cache build-base

WORKDIR /simple-file-server

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/bin/simple-file-server

COPY startup.sh /startup.sh

RUN chmod +x /startup.sh

VOLUME ["/data"]

ENTRYPOINT ["/startup.sh"]
