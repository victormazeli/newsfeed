# Newsfeed Backend


### How to run the project
1. Clone the project on to your local machine
2. Open the project using Goland or any Text editor with go support
3. Run the following command `go mod tidy` to download dependencies
4. Create an env file like so `development.env` for development purposes
5. Run `go run main.go` to start the project or `go run main.go -e development`
6. Navigate to `htt://0.0.0.0:PORT` to see api running
7. To deploy to production you create a `production.env` file and run `go run main.go -e production`,
this will start the server in production mode

### Project structure

````
newsfeedbackend
|
|--config
|    app.go
|    config.go
|    env.go
|
|--database
|       models
|           |--news.go
|           |--user.go
|
|       db.go
|
|--graph
|       generated
|           |--generated.go
|       model
|           |--models_gen.go
|       resolver.go
|       schema.graphqls
|       schema.resolvers.go
|
|--handlers
|      handler.go 
|
|--middlewares
|       auth.go
|       auth0.go
|       gqlcontext.go
|
|--redis
|       newscache.go
|
|--tests
|
|--utils
|       jwt.go
|       utils.go
|
|   main.go
|

````

### How to make request

1. Create new user

```
mutation{
  CreateNewUser(input:{email: "test@gmail.com", picture:"NA", full_name:"test name", user_id:"4VMDoJ2zAZ9fTPDh2a4nzs6In7OFU3eu@clients", topics:["sport","bitcoin"]}){
    _id
    email
    picture
    full_name
    topics{
      topic
    }
    user_id
    created_at
    updated_at
}
}

```
2. Get user info

```
query{
  GetUser{
    _id
    email
    picture
    full_name
    topics{
      topic
    }
    user_id
    created_at
    updated_at
  }
  
}
```
3. Get News feed

```
query{
  NewsFeed{
    author
    title
    description
    url
    urlToImage
    publishedAt
    content
    source{
      id
      name
    }
  }
  
}
```

Note: ensure to pass the token to the request header while making a request, Token are generated from auth0 a sample postman 
docs on generating a token. [click here](https://documenter.getpostman.com/view/12847208/2s935vjz6r)

```
{

Authorization: "YOUR_TOKEN"

}
```



### Project structure Break Down
1. `/config` contains all application configuration
2. `/database` contains database configuration and database model / schema
3. `/graph` contains all graphql related files and folder including graphql schema, resolvers etc.
4. `/handlers` contains methods and functions that performs specific actions for a particular use case
5. `/middlewares` contains all application middlewares
6. `/redis` contains redis configuration for the project
7. `/utils` contains common utility functions to be used all over the project
8. `main` this is the entry point of the application, includes route setup, server setup, configuration setup etc.

### Tools
1. Go (1.19)
2. Gin
3. Gqlgen
4. viper
5. redis
6. Mgm
7. Mongodb
8. resty

### References
1. Graphql (gqlgen): https://gqlgen.com/
2. Gin: https://gin-gonic.com/
3. mgm: https://github.com/Kamva/mgm
4. resty: https://github.com/go-resty/resty
5. viper: https://github.com/spf13/viper
