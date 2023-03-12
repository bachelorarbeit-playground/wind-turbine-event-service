# Generated Event Service

This is a generated event driven service designed to fit requirements of the scenario presented in the bachelor thesis 
**A Modelling Approach Focused on Event Processing to Generate Event Driven Applications** by Adrian Stefan Munteanu

![wind turbine service architecture](https://user-images.githubusercontent.com/23280777/224533389-cb7769d5-e292-476b-b214-e78f8bf1d67a.svg)

## Requirements
- Go
- Docker
- docker-compose
- [natscli](https://github.com/nats-io/natscli)

## Setup
After the requirements have been installed, start the NATS JetStream Cluster by running

```
docker-compose up
```

To install the Go project dependencies, run

```
go mod tidy
```

To get the event service running, run

```
go run cmd/main.go
```

## Scenario Simulation
In order to simulate the scenario, we must have the NATS JetStream Cluster and the event service running.
We must also open 2 terminal windows and subscribe to the output subjects. In each window run one the following commands: 

```
nats sub ingestionPipeline
```

```
nats sub anomalyDetection
```

Then, in a new terminal window, you can publish an input message to `rawWindData` with the `inputWindEvent.txt` content.

```
nats pub rawWindData "{\"park_id\":\"b86163c9-346e-45c9-93bb-dc77a22a5813\",\"turbine_id\":\"4deca05f-ceb7-474d-9c28-1e00b0c7521c\",\"region\":\"Berlin\",\"date\":\"2022-11-19\",\"interval\":19,\"timezone\":\"Europe/Berlin\",\"value\":0.03,\"availability\":95}"
```

This message should be forwarded to the `ingestionPipeline` subject and should show up in the subscription in one of the terminals opened at the beginning.

If we change the availablility property of the event to be `29` instead of `95`:

```
nats pub rawWindData "{\"park_id\":\"b86163c9-346e-45c9-93bb-dc77a22a5813\",\"turbine_id\":\"4deca05f-ceb7-474d-9c28-1e00b0c7521c\",\"region\":\"Berlin\",\"date\":\"2022-11-19\",\"interval\":19,\"timezone\":\"Europe/Berlin\",\"value\":0.03,\"availability\":29}"
```
it will show up in the `anomalyDetection` subject instead.
