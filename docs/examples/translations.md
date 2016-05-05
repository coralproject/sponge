# User Case : Washington Post

Washington Post uses Echo as a comment's platform. The import of data has been done through their web services. This is the translation of their data into the Coral schema.

## Users

```
Actor.id      => User.Source.user_id
Actor.title   => User.Name
Actor.status  => User.Status
Actor.avatar  => User.Avatar
Actor.markers => User.tags
Actor.roles   => User.roles
```

## Assets

```
Object.context.0.uri => Asset.Source.asset_id
Object.context.0.uri => Asset.URL
Object.published     => Asset.date_created
```

## Comments

```
Object.content      => Comment.Body
Object.status       => Comment.Status
Object.postedTime   => Comment.date_created
Object.Likes        => Comment.actions  - likes
Object.markers      => Comment.tags
Object.targets.0.id => Comment.Source.parent_id
Object.source       => Comment.Metadata.source
Object.markers      => Comment.Tags
Object.tags         => Comment.Metadata.Markers
```
