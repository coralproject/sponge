# Sponge

Move your existing community to Coral.

## Overview

Sponge is Coral's Data Import Layer.  It is designed to read data from a foreign _source_, translate the schema into coral conventions, and POST entities to [our service layer](https://github.com/coralproject/pillar) for inserting.

Where to get the data and how to translate it is expressed through _strategies_.

## Install

Read the [INSTALL file](https://github.com/coralproject/sponge/blob/master/INSTALL.md)

## Strategies

Strategies are json configuration files that contain all the information Sponge needs to get data from a source, translate it to the coral schema and send it to the service layer.

The strategy spec is still not complete. [Here's the living example](https://github.com/coralproject/sponge/blob/master/strategy/strategy.json.example)

## Code of conduct

Please be civil when discussing contributions to the Sponge code and the Coral Project. If in doubt, please consult our [Code of Conduct](https://the-coral-project.gitbooks.io/coral-bible/content/codeofconduct.html).

## Source support

Sponge currently only supports importing data from foreign databases.  For our plans to support other sources, [see the roadmap](ROADMAP.md).
