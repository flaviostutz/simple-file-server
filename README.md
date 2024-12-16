# simple-file-server

[<img src="https://img.shields.io/docker/pulls/flaviostutz/simple-file-server"/>](https://hub.docker.com/r/flaviostutz/simple-file-server)
[<img src="https://img.shields.io/docker/automated/flaviostutz/simple-file-server"/>](https://hub.docker.com/r/flaviostutz/simple-file-server)<br/>
[<img src="https://goreportcard.com/badge/github.com/flaviostutz/simple-file-server"/>](https://goreportcard.com/report/github.com/flaviostutz/simple-file-server)

Simple HTTP based file server.

Supports Uploading and Downloads.

For upload, if using POST Method, it will create a new file with a UUID name. If using PUT method, it will create the file with the same name as in request URI.

This server obeys ETag semantics. Check info at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag

## Usage

* Create docker-compose.yml

```yml
version: '3.7'

services:

  simple-file-server:
    build: .
    image: flaviostutz/simple-file-server
    ports:
      - "4000:4000"
    environment:
      - WRITE_SHARED_KEY=
      - READ_SHARED_KEY=
      - LOCATION_BASE_URL=http://localhost:4000
      - LOG_LEVEL=debug
```

* Run docker-compose build

* Run docker-compose up

* Execute

```bash
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X PUT http://localhost:4000/dir1/file1.json

curl http://localhost:4000/dir1/file1.json

curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST http://localhost:4000/dir2

```

## REST API

* POST /dir1/dir2 - creates a new file inside /dir1/dir2
  * Accepted headers:
    * "Mime-Type" - is stored so that later GET returns the same mime type on response
    * "X-Cache-Control" - is stored so that later GET returns "Cache-Control" header directive in response so that you can have caching proxies on front of this server
  * Returned headers:
    * "Location" - new generated file URL
    * "ETag" - file hash that you can be used later to verify if the file has changed or during file updates

* GET /dir1/dir2/c106fc67-eadb-4d91-beb1-66fc5c35e6a6 - get file contents along with ETag Header, Mime-Type and Cache-Control as set during creation/update

  * ETag - if you send Header "If-None-Match" the server will check if the current file's ETag hash is the same as it and if so, will return 304 Not Modified without contents, indicating that the client can still use the file contents in possession. If hash doesn't match, the new fresh contents will be returned

* PUT /dir1/dir2/c106fc67-eadb-4d91-beb1-66fc5c35e6a6 - updates and existing file or creates a new one with with a custom name (names are arbitrary and don't need to be an uuid)

  * ETag - when you send Header "If-Match" with an ETag hash Header, before updating any file contents, the server will check if this header value still matches current server file hash contents and fail if doesn't match. This is useful to guarantee that the file wasn't changed from the last time the application got it (other processes may have updated the file in the meanwhile and the application will know by handling this error)

## Automated tests

To run automated tests againts REST API, run
```
docker-compose -f docker-compose.test.yml up --build
```

This will build and run "simple-file-server" and another container with POSTMAN scripts to be executed against the server and check results.

For more info on how this "POSTMAN" runner works, go to https://github.com/flaviostutz/postman-runner

This tests will be run automatically if this repo is integrated to DockerHub with "Automated Tests" enabled (https://docs.docker.com/docker-hub/builds/automated-testing/).

## Generate Certificate (⚠ Only for local test purpose)

```bash
org=localhost
domain=localhost
openssl genpkey -algorithm RSA -out "$domain".key
openssl req -x509 -key "$domain".key -out "$domain".crt \
    -subj "/CN=$domain/O=$org" \
    -config <(cat /etc/ssl/openssl.cnf - <<END
[ x509_ext ]
basicConstraints = critical,CA:true
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
subjectAltName = IP:127.0.0.1
END
    ) -extensions x509_ext
```

*Note*:
subjectAltName IP:127.0.0.1 can not be replace by DNS:localhost as host verification will failled
cf : inline bool SSLClient::verify_host(X509 *server_cert) const
https://github.com/yhirose/cpp-httplib/blob/ae63b89cbf70481ae60515dfd95467e91eecd992/httplib.h#L9171C1-L9171C62

## Update root certificates (⚠ Only for local test purpose)

```bash
sudo trust anchor "$domain".crt
```

Or

```bash
sudo cp "$domain".crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```
