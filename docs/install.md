# Installation

## Run it from Docker Image


### Building image

To build the docker image run this command:

docker build -t "sponge:latest" -f Dockerfile ./

### Edit env.list

``STRATEGY_CONF=<path to strategy file>
PILLAR_URL=<url where pillar is running>``

`// DATABASE (optional if you want to overwrite strategy file values)
- DB_database= ""
- DB_username= ""
- DB_password= ""
- DB_host= ""
- DB_port= ""``

`// WEB SERVICE (optional if you want to overwrite strategy file values)
WS_appkey= ""
WS_endpoint= ""
WS_records= ""
WS_pagination= ""
WS_useragent= ""
WS_attributes= ""``

### Running the container

It will start imorting everything setup in the [strategy file](strategy.md).

``docker run --env-file env.list -d sponge``

## Development environment

You will need an instance of mysql running for the external source and an instance of [Pillar](http://github.com/coralproject/pillar) running to send the data to.

### Clone Repository

```
git clone git@github.com:CoralProject/sponge.git
```

### Configure


* Specify a strategy configuration file with the environment variable STRATEGY_CONF

```
export STRATEGY_CONF=/path/to/my/strategy.json
```

* Copy strategy file

There is one strategy file for publisher. It tells us how to do the transformation between the publisher's data and the coral data schema. It also tells us how to connect to the external publisher's source. There is a strategy file example in app/sponge/strategy.json.example.

```
cp app/sponge/strategy.json.example /path/to/my/strategy.json
```

* Specify the URL where pillar services are running

```
export PILLAR_URL=http://localhost:8080
```

### How to run

```
cd cmd/sponge
 go run main.go
```


### To build

```
cd cmd/sponge
go build
```

```
 ./sponge -h
```

will give you all the options to run it
