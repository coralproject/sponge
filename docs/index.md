Sponge is an [Extract, Transform and Load](https://en.wikipedia.org/wiki/Extract,_transform,_load) command line tool to get data from your community into the coral system. It does the translation needed to start using the Coral applications.

After [installing](/install.md) it, you will need to write the translations file, called [Strategy file](/readme.md). And have [Pillar](https://github.com/coralproject/pillar)

Setup the environment variables:

* STRATEGY_CONF with the path to the strategy file.
* PILLAR_URL with the URL where pillar is serving

## Packages included

* [Strategy](/strategy.md) reads the translations file.
* [Source](/source.md) does the extraction.
* [Fiddler](/fiddler.md) does the transformations.
* [Coral](/coral.md) send data to Pillar.
* [Sponge](/sponge.md) ties all the pieces together.

## Command line tool

### Usage:

  sponge --flag [command]

### Available Commands:

  * import      Import data to the coral database
  * index       Work with indexes in the coral database
  * show        Read the report on errors
  * version     Print the version number of Sponge
  * all         Import and Create Indexes

### Flags

```      --dbname="report.db": set the name for the db to read
      --filepath="report.db": set the file path for the report on errors (default is report.db)
  -h, --help[=false]: help for sponge
      --limit=9999999999: number of rows that we are going to import (default is 9999999999)
      --offset=0: offset for rows to import (default is 0)
      --onlyfails[=false]: import only the the records that failed in the last import(default is import all)
      --orderby="": order by field on the external source (default is not ordered)
      --query="": query on the external table (where condition on mysql, query document on mongodb). It only works with a specific --type. Example updated_date >'2003-12-31 01:02:03'
      --report[=false]: create report on records that fail importing (default is do not report)
      --type="": import or create indexes for only these types of data (default is everything)
```

## The Coral Documentation
