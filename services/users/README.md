## Users service

The users service provides an interface for creating and authenticating accounts for the application. A user is the central domain unit of the application - it defines authorization roles (user can only create posts for themselves for example), authentication (someone can only access the app if they have a user record and valid password for that user record), and acts as a reference for most/all data generated in the application (a user has many posts, has many followers, has followed many users, etc).

#### What is a user in the data model?

```
{
  username: text (pk),
  email: text (uniq),
  password: text,
  refresh_token: text
}
```

#### Access patterns

- Create a new user
- Retrieve or issue a user's access token via a "login" (username and password)
- TODO: Retrieve or issue a user's access token via a refresh token
