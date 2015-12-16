# Installation

### Clone Repository

```
git clone git@github.com:CoralProject/sponge.git
```

### Configure


* Copy strategy file

There is one strategy file for publisher. It tells us how to do the transformation between the publisher's data and the coral data schema. It also tells us how to connect to the external publisher's source. There is a strategy file example in app/sponge/strategy.json.example.

```
cp app/sponge/strategy.json.example app/sponge/strategy.json
```

### How to run

```
cd app/sponge
 go run main.go
```


### To build

```
cd app/sponge
go build
```

```
 ./sponge
```


### Example data to import

At ./scripts/LOADCSV.md we have the mysql commands on how to import the ./script/nyt_sample_data.csv example.
