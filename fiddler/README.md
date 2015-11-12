# Local Development Help Reference

## MongoDB

To manage mongodb we are using the shell and creating db, roles and user manually.

### To create a role

db.createRole(
  {
     role: "MyAdmin",
     privileges: [
       { resource: { db: "coral", collection: "" }, actions: [ "find", "update", "insert", "remove" ] }
     ],
     roles: []
  }
)


### To create a user

db.createUser({user:"username",pwd:"password",roles:[{role: "MyAdmin", db:"coral"}]})


#### Reference

* createRole: https://docs.mongodb.org/manual/reference/command/createRole/
