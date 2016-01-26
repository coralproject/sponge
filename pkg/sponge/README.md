## On how to import Data

## On how to create indexes

This feature create indexes on the coral database, depending on data in the strategy file.

### Setup endpoint for create index in Pillar


- Set CREATE_INDEX environment variable with the Pillar endpoint URI

```
export CREATE_INDEX=http://localhost:8080/api/import/index
```

### Add index's information to the strategy files

- Set the strategy file with the indexes that you want at the table's level.

```
"IndexBy": {
  "keys": {
    "field": "type of index"
  },
  "options": {} // set of options that control the creation of the index
},
```

More info at the (mongodb's create index definition)[https://docs.mongodb.org/manual/reference/method/db.collection.createIndex/#db.collection.createIndex].

### Run Sponge with Create Index Flags

- To create index on a specific collection

```
./sponge -index -type Asset
```

- To create indexes on all the collections

It create indexes only on the ones that you setup in the strategy file.

```
./sponge -index
```
