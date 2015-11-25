# Scripts

##### Fields in the NYT data sample

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

##### Creating the table.

```
CREATE TABLE nyt_comments (
  commentID INT NOT NULL,
  assetID INT NOT NULL,
  statusID INT,
  commentTitle TEXT,
  commentBody TEXT,
  userID INT,
  createDate DATE NOT NULL,
  updateDate DATE,
  approveDate DATE,
  commentExcerpt TEXT,
  editorsSelection VARCHAR(255),
  recommendationCount INT,
  replyCount INT,
  isReply INT,
  commentSequence VARCHAR(255),
  userDisplayName VARCHAR(255),
  userReply VARCHAR(255),
  userTitle VARCHAR(255),
  userLocation VARCHAR(255),
  showCommentExcerpt VARCHAR(255),
  hideRegisteredUserName INT,
  commentType VARCHAR(255),
  parentID VARCHAR(255),
  PRIMARY KEY (commentID)
)  ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

##### Loading the data

```
LOAD DATA INFILE 'nyt_sample_data.csv'
INTO TABLE nyt_comments
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;
```
