# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/goapp /app/
COPY en_punctuation.json /app/fixtures/
COPY en_stopwords.json /app/fixtures/


ENTRYPOINT ./goapp
