# snmpsim-restapi-go-client
[![Go Report Card](https://goreportcard.com/badge/github.com/inexio/snmpsim-restapi-go-client)](https://goreportcard.com/report/github.com/inexio/snmpsim-restapi-go-client)
[![GitHub license](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/inexio/check_eve_ng/blob/master/LICENSE)
[![GitHub code style](https://img.shields.io/badge/code%20style-uber--go-brightgreen)](https://github.com/uber-go/guide/blob/master/style.md)
[![GoDoc doc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/inexio/snmpsim-restapi-go-client)

## Description

Golang package - client library for the [snmpsim](https://github.com/etingof/snmpsim) [REST API](https://github.com/etingof/snmpsim-control-plane).

The library is used by [snmpsim-control-client](https://github.com/inexio/snmpsim-control-client) and [snmpsim-check](https://github.com/inexio/snmpsim-check).

This go client is an open-source library to communicate with the snmpsim REST API which is written in Python.

## Code Style

This project was written according to the **[uber-go](https://github.com/uber-go/guide/blob/master/style.md)** coding style.

## Features

### Management Client

- Create laboratory environments to test SNMP calls
- Configuration of laboratory environments
- Possibility to create Endpoints, Engines, Agents, Users and Tags
- Agents can be added to Laboratories
- Engines can be added to Agents
- Users and Endpoints can be added to Engines
- Tags can be applied to all of the above 
- Possibility to delete all objects linked to a tag (for cleanup purposes)

### Metrics Client

- Can check metrics of a lab environment
- Possibility to check processes, packet activity and message activity

## Requirements

The latest version of the snmpsim python module needs to be installed and configured.

Further information on how to download and configure snmpsim can be found [here](https://github.com/etingof/snmpsim).

Snmpsim-control-plane has to be installed and running. A guide for setting Control Plane can be found [here](http://snmplabs.com/snmpsim-control-plane/deployment.html).

To check if your setup works, follow the steps provided in the **'Tests'** section of this document.

## Installation

```
go get github.com/inexio/snmpsim-restapi-go-client
```

or 

``` 
git clone https://github.com/inexio/snmpsim-restapi-go-client.git
```

## Usage

### Management Client

```go
	//Create a new management api client
	client, err := snmpsimclient.NewManagementClient("https://127.0.0.1:8000")
	
	//Set http auth username and password (optional)
	err = client.SetUsernameAndPassword("httpAuthUsername", "httpAuthPassword")

	//Create a new lab
	lab, err := client.CreateLab("myLab") //optionally use CreateLabWithTag(..., tagId) [tagId as last param]

	//Create a new engine
	engine, err := client.CreateEngine("myEngine", "0102030405070809") //optionally use CreateEngineWithTag(..., tagId) [tagId as last param]

	//Create a new endpoint
	endpoint, err := client.CreateEndpoint("myEndpoint", "127.0.0.1", "1234") //optionally use CreateEndpointWithTag(..., tagId) [tagId as last param]

	//Create a new user
	user, err := client.CreateUser("uniqueUserIdentifier", "myUser", "", "", "", "") //optionally use CreateUserWithTag(..., tagId) [tagId as last param]

	//Add user to engine
	err = client.AddUserToEngine(engine.ID, user.ID)

	//Add endpoint to engine
	err = client.AddEndpointToEngine(engine.ID, endpoint.ID)

	//Create a new agent
	agent, err := client.CreateAgent("myAgent", "agent/data/dir") //optionally use CreateAgentWithTag(..., tagId) [tagId as last param]


	//Add engine to agent
	err = client.AddEngineToAgent(agent.ID, engine.ID)

	//Add agent to lab
	err = client.AddAgentToLab(lab.ID, agent.ID)

	//Set lab power on
	err = client.SetLabPower(lab.ID, true)
	
	//Delete lab
	err = client.DeleteLab(lab.ID)
```

### Metrics Client

```go
	//Create a new metrics api client
	client, err := snmpsimclient.NewMetricsClient("https://127.0.0.1:8001")

	//Set http auth username and password (optional)
	err = client.SetUsernameAndPassword("httpAuthUsername", "httpAuthPassword")

	//Get all packet metrics
	packets, err := client.GetPackets(nil)

	//Get all possible filters for packets
	packetFilters, err := client.GetPacketFilters()
	
	//Get all packets for endpoint "127.0.0.1:1234"
	filters := make(map[string]string)
	filters["local_address"] = "127.0.0.1:1234"
	packets, err := client.GetPackets(filters)

	//Get all message metrics
	messages, err := client.GetMessages(nil)
```



### Tests

Our library provides a few unit and integration tests. To use these tests, the yaml config files in the test-data directory must be adapted to your setup.

In order to run these test, run the follwing command inside root directory of this repository:

```
go test
```



If you want to check if your setup works, run:

```
go test -run TestManagementClient_buildUpSetupAndTestIt
```



## Getting Help

If there are any problems or something does not work as intended, open an issue on GitHub.

## Contribution

Contributions to the project are welcome.

We are looking forward to your bug reports, suggestions and fixes.

If you want to make any contributions make sure your go report does match up with our projects score of **A+**.

When you contribute make sure your code is conform to the **uber-go** coding style.

Happy Coding!



