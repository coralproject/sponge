# Sponge

import "github.com/coralproject/sponge/pkg/sponge"

Imports external source database into local source, transform it and send it to the coral system (called pillar).

## constants

``const (
	VersionNumber = 0.1
)
``

Updated version of Sponge


## variables

``var (
	dbsource source.Sourcer
	uuid     string
	options  source.Options
)
``

## func AddOptions

``func AddOptions(limit int, offset int, orderby string, query string, types string, importonlyfailed bool, reportOnFailedRecords bool, reportdbfile string)``

Sets options for how sponge is going to be running. The options are:

	*	Limit: limit for the query
	*	Offset: offset for the query
	*	Orderby:  order by this field
	*	Query:  we use this field if we want to specific a filter on WHERE for mysql/postgresql and Find for MongoDB
	*	Types: it specifies which entities to import (default is everything)
	*	Importonlyfailed: import only the documents that are in the report
	*	ReportOnFailedRecords: create a report with all the documents that failed the import
	*	Reportdbfile: name of the file for the report on documents that fail the import



## func Import

  ``func Import()``

Gets data, transform it and send it to pillar. It based everything on STRATEGY_CONF's environment variable and PILLAR_URL environment variable.

## func CreateIndex

  ``func CreateIndex(collection string)``

Create index on the collection 'collection'. This feature create indexes on the coral database, depending on data in the strategy file.

For example:

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

More info at the [mongodb's create index definition](https://docs.mongodb.org/manual/reference/method/db.collection.createIndex/#db.collection.createIndex).
