# Installation

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
