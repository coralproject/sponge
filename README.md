# Sponge

Get your existing community to Coral. This is a command line tool to import comments, authors, assets and other entities into the coral system.

## Overview

Sponge is Coral's Data Import Layer.  It is designed to read data from a foreign _source_, translate the schema into coral conventions, and POST entities to [our service layer](https://github.com/coralproject/pillar) for inserting.

Where to get the data and how to translate it is expressed through _strategies_.

## Install

Clone this repository , install GO and build main.go. More information at the [INSTALL's file](https://github.com/coralproject/sponge/blob/master/docs/INSTALL.md).

## Documentation

We use [Sponge's Read The Docs](https://sponge.readthedocs.io) to publish documentation for this repository.

## Dependencies

You should be vendoring the packages you choose to use. We recommend using [govendor](https://github.com/kardianos/govendor). This tool will vendor from the vendor folder associated with this project repo for the dependencies in use. It is recommended to use a project based repo with a single vendor folder for all your dependencies.

## Examples

We distribute [examples on strategy file](/examples/)  as well as [documentation on user's cases](/docs/examples) on this repository.

## Roadmap

Sponge currently only supports importing data from mysql, postgresql, mongodb or web services.  For our plans to support other sources, [see the roadmap](ROADMAP.md).

## Code of conduct

Please be civil when discussing contributions to the Sponge code and the Coral Project. If in doubt, please consult our [Code of Conduct](https://the-coral-project.gitbooks.io/coral-bible/content/codeofconduct.html).

## License

Copyright 2016, The Coral Project. [MIT License](/LICENSE)

## Mantainer

Sponge is maintained by The Coral Project. Check the [Contributions File](/CONTRIBUTORS.md) for more information.
