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

    curl -XPOST -H"Content-Type:application/json;charset=utf-8" @sample-native-article.json http://localhost:8080/map

## Build and deployment

* Built by Docker Hub on merge to master: [coco/methode-article-image-set-mapper](https://hub.docker.com/r/coco/methode-article-image-set-mapper/)
* CI provided by CircleCI: [methode-article-image-set-mapper](https://circleci.com/gh/Financial-Times/methode-article-image-set-mapper)

## Service endpoints

### POST

    curl -XPOST @sample-native-article.json http://localhost:8080/map

The expected response will be an image-set.

## Admin endpoints:

* `/__gtg`
* `/__health`
* `/__build-info`
* `/__ping`

Healthchecks check that the app can read from a kafka topic and write to another.
