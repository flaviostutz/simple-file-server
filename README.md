# simple-file-server

Simple HTTP based file server.

Supports Uploading and Downloads.

For upload, if using POST Method, it will create a new file with a UUID name. If using PUT method, it will create the file with the same name as in request URI.

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
      - LOG_LEVEL=debug
```

* Run docker-compose up

* Execute

```bash
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST http://localhost:4000/dir1/file1.json

curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X PUT http://localhost:4000/dir2
```
