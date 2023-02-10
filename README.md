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
