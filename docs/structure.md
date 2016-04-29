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
