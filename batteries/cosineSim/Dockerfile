FROM golang:alpine AS build-env

ADD . /src
COPY en_punctuation.json /src/fixtures/
COPY en_stopwords.json /src/fixtures/

# compile app
RUN cd /src && go build -o goapp

# COPY /src/goapp /app/

ENTRYPOINT ./src/goapp