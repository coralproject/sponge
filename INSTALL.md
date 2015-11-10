# Installation

### Clone Repository

```
git clone git@github.com:CoralProject/sponge.git
```

### Configure


* Copy configuration file

```
cp config/config.json.example config/config.json
```

* Configure database

	* Name: name of the Publisher you are configuring for.
	* Strategy:  type of database and tables you are importing.
	* Credentials: Array of credentials with information about each database you are importing from and to. Adapter is the kind of DBMS and type is source or local.

* To add a new collection/table

	* Modify config.json with collection: table in the strategy.Tables
	* Add new model to models interface.  This needs to be done automatically.
	* Adds new collection to the utils New func. This needs to be done automatically.

### How to run

```
 go run main.go
```


### To build

```
 go build -o import
```

```
 ./import
```


### Example data to import

At ./scripts/LOADCSV.md we have the mysql commands on how to import the ./script/nyt_sample_data.csv example.
