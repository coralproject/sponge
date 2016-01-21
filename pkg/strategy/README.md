# Strategy configuration

This is a json file that describes the transformations that need to happen to the external data to get into the coral system.

## Name

The name of the strategy that we are describing.

## Map

Contains all the information to map fields from the external source to our local coral database.

### Foreign

It describes the foreign database source. Right now we have only "mysql" available.

### DateTimeFormat

If your source have date time fields, we need to know how to parse them. You should write the representation of 2006 Mon Jan 2 15:04:05 in the desired format. More more info read [pkg-constants](https://golang.org/pkg/time/#pkg-constants)

### Tables

It describe all the tables and its transformations.

#### Name of the table

Example:
```
"asset": {
  "Foreign": "crnr_asset",
  "Local": "asset",
  "Priority": 1,
  "OrderBy": "assetid",
  "ID": "assetid",
  "IndexBy": {
    "keys": {
      "asseturl": "text"
    },
    "options": {}
  },
  "Fields": [
    {
      "foreign": "assetid",
      "local": "asset_id",
      "relation": "Source",
      "type": "int"
    },
    {
      "foreign": "asseturl",
      "local": "url",
      "relation": "Identity",
      "type": "int"
    },
    {
      "foreign": "updatedate",
      "local": "date_updated",
      "relation": "ParseTimeDate",
      "type": "dateTime"
    },
    {
      "foreign": "createdate",
      "local": "date_created",
      "relation": "ParseTimeDate",
      "type": "dateTime"
    }
  ],
  "Endpoint": "http://localhost:8080/api/import/asset"
}
```

##### Foreign

The name of the foreign table or collection.

##### Local

The collection to where we are importing this table into.

##### Priority

This is a number that specifies which table to import first.

##### OrderBy

A default order by when quering the foreign source.

##### ID

The identifier field for the foreign table. We use this field when we need to import only some records and not the whole table.

##### Endpoint

This is the endpoint in the coral system were we are going to push the data for this table into.

##### Fields

All the fields for the table with the mapping

###### Foreign

The name of the field in the foreign table.

###### Local

The name of the field in our local database.

###### Relation

The relationship between the foreign field and the local one. We have this options:
- Identity: when the value is the same
- Source: when it needs to be added to our source struct for the local table (the orignal identifiers have to go into source)
- ParseTimeDate: when we need to parse the foreign value as date time.

###### Type

The type of the value we are converting.

## Credentials

This has the credentials for the foreign database to pull data from.

### adapter

It tells us which drivers we need to use to pull data. Right now we have the "mysql" option.

### type

Right now it is always "foreign" but it could tell us which type of credential this one is.
