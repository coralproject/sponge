## On how to import Data

## On how to create indexes

This feature create indexes on the coral database, depending on data in the strategy file.

### Setup endpoint for create index in Pillar


- Have PILLAR_URL environment variable set with the Pillar server url

```
export PILLAR_URL=http://localhost:8080
```

### Add index's information to the strategy files

- Set the strategy file with the indexes that you want at the table's level.

```
"Index": [
  {
    "name": "asset-url",
    "keys": ["asseturl"],
    "unique": "true",
    "dropdups": "true"
  }
],
```

More info at the (mongodb's create index definition)[https://docs.mongodb.org/manual/reference/method/db.collection.createIndex/#db.collection.createIndex].
