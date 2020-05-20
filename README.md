# simple-file-server

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

* Run docker-compose up

* Execute

```bash
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X PUT http://localhost:4000/dir1/file1.json

curl http://localhost:4000/dir1/file1.json

curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST http://localhost:4000/dir2

```

## REST API

* POST /dir1/dir2 - creates a new file inside /dir1/dir2 and returns the generated file URL in HTTP Header "Location" along with Header ETag, that you can use later to verify if the file has changed or during file updates. Mime type is stored so that GET returns the same mime type.

* GET /dir1/dir2/c106fc67-eadb-4d91-beb1-66fc5c35e6a6 - get file contents along with ETag Header and mime type set during creation/update

  * ETag - if you send Header "If-None-Match" the server will check if the current file's ETag hash is the same as it and if so, will return 304 Not Modified without contents, indicating that the client can still use the file contents in possession. If don't match, the new contents will be returned

* PUT /dir1/dir2/c106fc67-eadb-4d91-beb1-66fc5c35e6a6 - updates and existing file or creates a new one with with a custom name (names are arbitrary and don't need to be an uuid)

  * ETag - if you send Header "If-Match" with an ETag hash value, if the file exists and will be updated, it will be checked if the current file on server still has this ETag set and the call will fail if not. This is useful to guarantee that the file wasn't changed from the last time the application got it (other processes may have updated the file in the meanwhile and the application will know by this error)
