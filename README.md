# Sponge

[Sponge](https://github.com/coralproject/sponge) is Coral's data import layer. It is an Extract, Transform, and Load command line tool designed to:

* Read data from a foreign source,
* Translate the schema into Coral conventions, and
* POST entities to our service layer for insertion.

All of the [Sponge documentation](https://coralprojectdocs.herokuapp.com/sponge/) (including installation instructions) can be found in the [Coral Project Documentation](https://coralprojectdocs.herokuapp.com/).

<<<<<<< HEAD
Where to get the data and how to translate it is expressed through _strategies_.

## Install

### Build from Source

#### Install Go

To build Sponge you will need [Go 1.6.1+](https://golang.org/dl/).

First download `sponge` without installing:

```
$ go get -d github.com/coralproject/sponge

$ cd $GOPATH/src/github.com/coralproject/sponge
```

Then install `sponge` and its dependencies:

```
$ make install
```

More information at the [INSTALL's file](https://github.com/coralproject/sponge/blob/master/docs/INSTALL.md).

## Documentation

We use [Sponge's Read The Docs](https://sponge.readthedocs.io) to publish documentation for this repository.

## Dependencies

You should be vendoring the packages you choose to use. We recommend using [govendor](https://github.com/kardianos/govendor). This tool will vendor from the vendor folder associated with this project repo for the dependencies in use. It is recommended to use a project based repo with a single vendor folder for all your dependencies.

## Examples

We distribute [examples on strategy file](/examples/)  as well as [documentation on user's cases](/docs/examples) on this repository.

## Roadmap

Sponge currently only supports importing data from mysql, postgresql, mongodb or web services.  For our plans to support other sources, [see the roadmap](ROADMAP.md).

## Code of conduct

Please be civil when discussing contributions to the Sponge's code source and the Coral Project. If in doubt, please consult our [Code of Conduct](https://the-coral-project.gitbooks.io/coral-bible/content/codeofconduct.html).

## License

Copyright 2016, The Coral Project. [MIT License](/LICENSE)

## Mantainer

Sponge is maintained by The Coral Project. Check the [Contributions](/CONTRIBUTORS.md) and [Authors](/AUTHORS.md) files for more information.
=======
The Sponge documentation [lives in Github](https://github.com/coralproject/docs/tree/master/docs_dir/sponge) in the `coralproject/docs/docs_dir/sponge` repository.
>>>>>>> 4cb1a797eaa6fc92d9dfe07769dd71e65d0bfbc6
