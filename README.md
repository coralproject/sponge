# Sponge

Move your existing community to Coral.

## Overview

Sponge is Coral's Data Import Layer.  It is designed to read data from a foreign _source_, translate the schema into coral conventions, and POST entities to [our service layer](https://github.com/coralproject/pillar) for inserting.

Where to get the data and how to translate it is expressed through _strategies_.

## Install

Read the [INSTALL file](https://github.com/coralproject/sponge/blob/master/INSTALL.md)

## Docker Image

To build the docker image run this command:

docker build -t "sponge:latest" -f Dockerfile ./

To run it

docker run --env-file env.list -d sponge

##### Folder structure

```
.
+-- cmd
|  +-- sponge
|       +-- cmd
+-- data                 -> examples for strategy files
+-- pkg
|  +-- coral             -> package to send data to pillar's endpoints
|  +-- fiddler           -> it does all the transformation needed for the data
|  +-- report            -> creates (and read) a report file with failed records
|  +-- source            -> drives to import data from different databases
|  +-- sponge            -> it imports data, transform it and send it to pillar
|  +-- strategy          -> parse the strategy file
|  +-- webservice        -> utility that encapsulate the requests to a web service
+-- tests                -> fixtures to use in the tests
+-- vendor               -> vendoring of all the external packages needed
```


## Strategies

Strategies are json configuration files that contain all the information Sponge needs to get data from a source, translate it to the coral schema and send it to the service layer.

The strategy spec is still not complete. You can read this [living example](https://github.com/coralproject/sponge/blob/master/strategy/strategy_mysql.json.example) or [more information on Strategies](https://github.com/coralproject/sponge/tree/master/pkg/strategy).

## About LOGGING

We are using (Ardanlabs Log's package)[https://github.com/ardanlabs/kit/tree/master/log] for all the tools we are developing in GO.

#### Spec:

###### Logging levels:

    Dev: to be outputted in development environment only
    User (prod): to be outputted in dev and production environments

###### All logs should contain:

    context uuid to identify a particular execution (aka, run of Sponge or a Request/Response execution from a web server.)
    the name of the function that is executing
    a readable message including relevant data

Logs should write to stdout so they can be directed flexibly.

More information the [reef wiki](https://github.com/coralproject/reef/wiki/Application-Logging).

#### Development

We welcome community contributions. If you're thinking about making more than a minor change, check in with the Coral team via Github issues to avoid programming conflicts.

We will not accept commits or pushes to the `master` branch, as the latest version of master is automatically deployed. Any direct push to master will be reverted.

Sequester all work in pull requests

  1. create a new branch with `git checkout -b your-fancy-branch-name`
  2. make a trivial change, and commit back to your branch with `git add ./your-changed-file.js` and `git commit -m "a commit message here"`
  3. push your changes to github with `git push origin your-fancy-branch-name`
  4. on github.com, you should see a button to create a pull request from your new branch
  5. There will be public code reviews before we merge any PRs into master


## Code of conduct

Please be civil when discussing contributions to the Sponge code and the Coral Project. If in doubt, please consult our [Code of Conduct](https://the-coral-project.gitbooks.io/coral-bible/content/codeofconduct.html).

## Source support

Sponge currently only supports importing data from foreign databases.  For our plans to support other sources, [see the roadmap](ROADMAP.md).
