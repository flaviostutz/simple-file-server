FROM node:13.10.1-alpine3.11

ENV MAX_FILESIZE_KB '10240'
ENV UPLOAD_PATH '/files'

RUN npm install -g upload-test-server

ADD startup.sh /

CMD [ "/startup.sh" ]

