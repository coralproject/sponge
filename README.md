# Data Import Layer

The import layer will contain strategies for pulling data from existing systems into Coral.  The data will be pulled into a mongo db without transformation initially.


### NY Times Strategy

type: user, item, boolAction, associationAction, externalReference

Schema
```{
  database: nyt,

  "comment": {
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

          relation: 'self',
          model: Comment
        },
        {
          name: statusID,
          type: int,
          local: ModerationStatus,

          relation: 'self',
          model: Comment
        },
        {
          name: commentBody
          type: []byte,
          local: Content,

          relation: 'self',
          model: Comment
        },
        {
          name: createDate,
          type: TimeDate,
          local: sourceCreateDate,

          relation: 'self',
          model: Comment
        },
        {
          name: updateDate,
          type: TimeDate,
          local: sourceUpdateDate,

          relation: 'self',
          model: Comment
        },
        {
          name: approveDate,
          type: TimeDate,
          local: sourceApproveDate,

          relation: 'self',
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

          relation: 'belongsTo',
          model: User
        },

        {
          name: AssetId,
          type: []byte,
          local: AssetId,

          relation: 'belongsTo',
          model: Asset
        }.
      ],    
    },

  "user": {
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

        relation: 'self',
        model: User
      },
      {
        name: UserDisplayName,
        type: string
        local: Name,

        relation: 'self',
        model: User
      },
      {
        name: UserLocation,
        type: string
        local: Location,

        relation: 'self',
        model: User
      }
    ]
  }
}```

#### NYT Comment's fields

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
