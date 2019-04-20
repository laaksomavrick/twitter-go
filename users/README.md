## Users

The users logical service provides a set of functions for creating and authenticating accounts for the application. A user is the central domain unit of the application - it defines authorization roles (user can only create posts for themselves for example), authentication (someone can only access the app if they have a user record and valid password for that user record), and acts as a reference for most/all data generated in the application (a user has many posts, has many followers, has followed many users, etc).

#### What is a user in the data model?

```
{
  username: pk, hash
  email: s (uniq),
  password: s (bcrypted),
  refresh_token: s
}
```

#### Access patterns

- Create a new user
- Retrieve or issue a user's access token via a "login" (username and password) _todo_
  - this should generate a new refresh token
- Retrieve or issue a user's access token via a refresh token _todo_
- Verify an access token against the one stored in the database _todo_
- Deactivate / reactivate a user _todo_

#### Dev flow

`make`
`serverless deploy`
`serverless invoke -f $fnname`
`serverless remove -v`
