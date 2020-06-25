# UPP - Splunk Event Reader

Maps inline image-sets from bodies of Methode articles.

## Primary URL

<https://upp-prod-delivery-glb.upp.ft.com/__methode-article-image-set-mapper/>

## Service Tier

Platinum

## Lifecycle Stage

Production

## Delivered By

content

## Supported By

content

## Known About By

- dimitar.terziev
- hristo.georgiev
- elina.kaneva
- tsvetan.dimitrov
- robert.marinov
- georgi.ivanov

## Host Platform

AWS

## Architecture

This service reads messages from the NativeCmsPublicationEvents kafka topic and maps them to image sets, then writes them to the CmsPublicationEvents topic. It also exposes an HTTP endoint for testing the mapping functionality.

## Contains Personal Data

No

## Contains Sensitive Data

No

## Dependencies

- upp-kafka

## Failover Architecture Type

ActiveActive

## Failover Process Type

FullyAutomated

## Failback Process Type

FullyAutomated

## Failover Details

The service is deployed in both Delivery clusters. The failover guide for the cluster is located here:
<https://github.com/Financial-Times/upp-docs/tree/master/failover-guides/delivery-cluster>

## Data Recovery Process Type

NotApplicable

## Data Recovery Details

The service does not store data, so it does not require any data recovery steps.

## Release Process Type

PartiallyAutomated

## Rollback Process Type

Manual

## Release Details

Failover is not needed during release.

## Key Management Process Type

Manual

## Key Management Details

To access the service clients need to provide basic auth credentials.
To rotate credentials you need to login to a particular cluster and update varnish-auth secrets.

## Monitoring

Service in UPP K8S delivery clusters:

- Delivery-Prod-EU health: <https://upp-prod-delivery-eu.upp.ft.com/__health/__pods-health?service-name=methode-article-image-set-mapper>
- Delivery-Prod-US health: <https://upp-prod-delivery-us.upp.ft.com/__health/__pods-health?service-name=methode-article-image-set-mapper>

## First Line Troubleshooting

<https://github.com/Financial-Times/upp-docs/tree/master/guides/ops/first-line-troubleshooting>

## Second Line Troubleshooting

Please refer to the GitHub repository README for troubleshooting information.
