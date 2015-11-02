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
The main loop that keeps the data up to date.  Gets _slices_ of records based on _updated_at_ timestamps to account for changing records.

For each table or api in the strategy:

* Ensure maximum rate limit is not met
* Determine which slice of data to get next
	* Find the last updated timestamp in the _log collection_ (or the collection itself?)
* Use the strategy to request the slice (either db query or api call)
	* Update the rate limit counter
* For each record returned
	* Check to ensure the document isn't already added
	* If not, add the document and kick off _import actions_
	* If it's there, update the document
	* Update the _log collection_


### Application States

#### Initialization Phase

* Initialize the _log_ collection with time that initialization began
* Kick off _synchronization loop_ at _initialization frequency_

#### Synchronization Phase

* Kick off the _synchronization loop_ at the _synchronization frequency_

### Import Actions

Actions performed on each document inserted

* Insert node(s) vertice(s) into neo4j
* Run data through data science pipelines

### Rate Limit Counter

A routine that keep as sliding count of how many requests were made in the past time frame based on the strategy.  

Exposes isOkToQuery() to determine if we are currently at the limit.

Each request sends a message to this routine each time a request is made.


### NY Times Strategy

type: user, item, boolAction, associationAction, externalReference

#### Schema

```
{
  database: nyt,
  comment: {
  	type: item,
    table: {
        foreign: comments,
        local: comments
      },
      fields: [
        {
          name: CommentID,
          type: int,
          primaryKey: true,
          local: id,
          relation: ,
          model: Comment
        },
        {
          name: statusID,
          type: int,
          local: ModerationStatus,
          relation: '',
          model: Comment
        },
        {
          name: commentBody
          type: []byte,
          local: Content,
          relation: '',
          model: Comment
        },
        {
          name: createDate,
          type: TimeDate,
          local: sourceCreateDate,
          relation: '',
          model: Comment
        },
        {
          name: updateDate,
          type: TimeDate,
          local: sourceUpdateDate,
          relation: '',
          model: Comment
        },
        {
          name: approveDate,
          type: TimeDate,
          local: sourceApproveDate,
          relation: '',
          model: Comment
        },
        {
          name: recommendationCount,
          type: int
          ????
        },
        {
          name: ParentID,
          type: int,
          local: ParentId,
          relation: 'hasMany',
          model: Comment
        },
        {
          name: UserID,
          type: int,
          local: id,
          relation: belongsTo,
          model: User
        },
        {
          name: AssetId,
          type: []byte,
          local: AssetId,
          relation: belongsTo,
          model: Asset
        }.
      ],    
    },
  user: {
    type: item,
    table: {
      foreign: comments,
      local: users
    },
    fields: [
      {
        name: UserID,
        type: int,
        primaryKey: true,
        local: id,
        relation: '',
        model: User
      },
      {
        name: UserDisplayName,
        type: string
        local: Name,
        relation: '',
        model: User
      },
      {
        name: UserLocation,
        type: string
        local: Location,
        relation: '',
        model: User
      }
    ]
  }
}
```

#### NYT Comment's Table's Fields

1. commentID: comment's identifier. Numeric. - YES
2. assetID: Numeric. Article's identifier. - YES
3. statusID: Numeric. Options: [2, 3]  - moderated(1), rejected(2), flagged(3) - YES (ModerationStatus)
4. commentTitle: String. Maybe NULL  - deprecated
5. commentBody: String. Maybe NULL - YES (Content)
6. userID: user that commented. Numeric.  (registration id) - YES
7. createDate: TimeDate   - YES (sourceCreatedDate)
8. updateDate: TimeDate   - YES
9. approveDate: TimeDate  - YES
10. commentExcerpt: Summary? String. Maybe NULL  - deprectaed
11. editorsSelection: Boolean  - nyt picks - YES (this is a BooleanAction)
12. recommendationCount: Numeric  - how many people recommended the comment  - YES
13. replyCount: Numeric  - how many people reply to that comment - NO
14. isReply: Boolean  - if it is a reply or not - NO
15. commentSequence: Numeric  - NO
16. userDisplayName: The user's name. String.  Deprecated. - YES (into User model)
17. userURL: Deprecated.
18. userTitle:  Deprecated.
19. userLocation: City. String. Deprecated. - YES (into User model - json list.. )
20. showCommentExcerpt: Boolean  - NO
21. hideRegisteredUserName: Boolean  - Deprecated?? - NO (check if any is not zero)
22. commentType: Deprecated. - NO
23. parentID: The id of the parent's comment. - YES
