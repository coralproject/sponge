# Data Import Layer

The import layer tries to keep our _local_ data in sync with _foreign_ data.

It will contain strategies for pulling data from existing systems into Coral.  The data will be pulled into a mongo db without transformation.

## Overview

This program imports data from a _foreign_ source into a local mongo database.  _Strategies_ are used to describe how to access the data from the _foreign_ source.  Data is stored in mongo collections verbatim.

### Foreign Source Types
*DB connection* - Databases that we can directly connect to.

* MySQL

*APIs* - Http(s) endpoints that provide the data.

* Disqus - https://disqus.com/api/docs/
* Wordpress Core - https://codex.wordpress.org/XML-RPC_WordPress_API/Comments
* Lyvewire - http://answers.livefyre.com/developers/api-reference/
* Facebook - https://developers.facebook.com/docs/graph-api/reference/v2.5/comment
  
### Strategies

* Field translation table: map of 
* Data relationships
* Limit of rows to request at a time
* Update Frequency:  
* Maximum Request Limits

*DB Connection* 

* Connection information: host, port, username, password, database
* Tables to import: list of tables that will be duplicated in our local collection

*APIs*

(To be clarified as we analyze different APIs)
* Endpoints to hit
* Request parameters to achieve the slices


### Synchronization Loop
The main loop that keeps the data up to date.  Each time through the loop:

* Look at the _log collection_ to determine which slice of data to request next
* Use the strategy to request the slice (either db query or api call)
* For each record returned
	* Check to ensure the document isn't already added 
	* If not, add the document
	* If it's there, log a sync error
* Update the _log collection_ 

### Application States

#### Initialization Phase 

* Initialize the _log_ collection with time that initialization began
* Kick off _synchronization loop_ at _initialization frequency_

#### Synchronization Phase 

* Kick off the _synchronization loop_ at the _synchronization frequency_
