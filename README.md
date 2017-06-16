# methode-article-image-set-mapper

[![Circle CI](https://circleci.com/gh/Financial-Times/methode-article-image-set-mapper/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/methode-article-image-set-mapper/tree/master)[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/methode-article-image-set-mapper)](https://goreportcard.com/report/github.com/Financial-Times/methode-article-image-set-mapper) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/methode-article-image-set-mapper/badge.svg)](https://coveralls.io/github/Financial-Times/methode-article-image-set-mapper)

## Introduction

Maps inline image-sets from bodies of Methode articles.

## Installation

Download the source code, dependencies and test dependencies:

    go get -u github.com/kardianos/govendor
    go get -u github.com/Financial-Times/methode-article-image-set-mapper
    cd $GOPATH/src/github.com/Financial-Times/methode-article-image-set-mapper
    govendor sync
    govendor test -v -race
    go build

## Run

    ./methode-article-image-set-mapper

Options:

    --app-system-code="methode-article-image-set-mapper"  System Code of the application ($APP_SYSTEM_CODE)
    --app-name="methode-article-image-set-mapper"         Application name ($APP_NAME)
    --port="8080"                                         Port to listen on ($APP_PORT)
    --queue-addresses="http://ip-172-24-74-51.eu-west-1.compute.internal:8080"  Hostname and port to connect to kafka proxy at ($Q_ADDR)
    --group="methode-article-image-set-mapper"            Kafka consumer group to use for incoming messages (Q_GROUP)
    --read-topic="NativeCmsPublicationEvents"             Topic to read messages that need mapping ($Q_WRITE_TOPIC)
    --read-queue="kafka"                                  HTTP host header used for routing in UP stack to reach the kafka container ($Q_READ_QUEUE)
    --write-topic="CmsPublicationEvents"                  Topic to write mapped messages to ($Q_WRITE_TOPIC)
    --write-queue="kafka"                                 Same as for read-queue ($Q_WRITE_QUEUE)

## Try:

    ssh -L 8083:localhost:8080 core@rj-tunnel-up.ft.com

    export Q_ADDR="http://localhost:8083"
    export Q_GROUP="methode-article-image-set-mapper-local"
    export Q_READ_TOPIC=NativeCmsPublicationEvents
    export Q_READ_QUEUE=kafka
    export Q_WRITE_TOPIC=CmsPublicationEvents
    export Q_WRITE_QUEUE=kafka

    curl -XPOST -H"Content-Type:application/json;charset=utf-8" -H"X-Request-Id:tid_test" -d @sample-methode-native-article-c17e8abe-1df8-11e7-942c-4a4c42b3072e.json http://localhost:8080/map

## Build and deployment

* Built by Docker Hub on merge to master: [coco/methode-article-image-set-mapper](https://hub.docker.com/r/coco/methode-article-image-set-mapper/)
* CI provided by CircleCI: [methode-article-image-set-mapper](https://circleci.com/gh/Financial-Times/methode-article-image-set-mapper)

## Service endpoints

### /map

### POST

Request:

The input should be a methode article in native format, as it comes from methode directly or similar to what is stored in the native-store.

    curl -XPOST -H"Content-Type:application/json;charset=utf-8" -H"X-Request-Id:tid_test" -d @sample-methode-native-article-c17e8abe-1df8-11e7-942c-4a4c42b3072e.json http://localhost:8080/map

Response:

The expected response will be an array of image-sets.

```
HTTP/1.1 200 OK
Content-Type: application/json;charset=utf-8
Date: Mon, 22 May 2017 15:34:15 GMT
Content-Length: 1055

[
  {
    "uuid": "4ec94836-0d00-325d-9005-c9aa67f68963",
    "identifiers": [
      {
        "authority": "http://api.ft.com/system/FTCOM-METHODE",
        "identifierValue": "4ec94836-0d00-325d-9005-c9aa67f68963"
      }
    ],
    "members": [
      {
        "uuid": "41614f4c-13c5-11e7-9469-afea892e4de3"
      },
      {
        "uuid": "4258f26a-13c5-11e7-9469-afea892e4de3"
      },
      {
        "uuid": "3ff3b7a8-13c5-11e7-9469-afea892e4de3"
      }
    ],
    "publishReference": "tid_test",
    "lastModified": "2017-05-22T02:59:39.195Z",
    "publishedDate": "2017-04-10T03:29:14.000Z",
    "firstPublishedDate": "2017-04-10T03:29:14.000Z",
    "canBeDistributed": "yes"
  },
  {
    "uuid": "8c07916c-2577-37b6-b477-291094f992ee",
    "identifiers": [
      {
        "authority": "http://api.ft.com/system/FTCOM-METHODE",
        "identifierValue": "8c07916c-2577-37b6-b477-291094f992ee"
      }
    ],
    "members": [
      {
        "uuid": "41614f4c-13c5-11e7-9469-afea892e4de3"
      },
      {
        "uuid": "4258f26a-13c5-11e7-9469-afea892e4de3"
      },
      {
        "uuid": "3ff3b7a8-13c5-11e7-9469-afea892e4de3"
      }
    ],
    "publishReference": "tid_test",
    "lastModified": "2017-05-22T02:59:39.195Z",
    "publishedDate": "2017-04-10T03:29:14.000Z",
    "firstPublishedDate": "2017-04-10T03:29:14.000Z",
    "canBeDistributed": "yes"
  }
]
```

## Admin endpoints:

* `/__gtg`
* `/__health`
* `/__build-info`
* `/__ping`

Healthchecks check that the app can read from a kafka topic and write to another.
