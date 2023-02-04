package main

import (
	"flag"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"newsfeedbackend/config"
	"newsfeedbackend/database"
	"newsfeedbackend/graph"
	"newsfeedbackend/graph/generated"
	"newsfeedbackend/middlewares"
	"newsfeedbackend/redis"
)

/*
This is the entry file to the application, here we created the graphql server
*/

const defaultPort = "8080"

// Defining the Graphql handler
func graphqlHandler(env *config.Env) gin.HandlerFunc {
	// pass the env to the resolver since the resolver instance is where we can initialize dependencies
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Env: env,
	}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	// toggle between different environment using command line
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
	}
	flag.Parse()
	//
	app := config.App(*environment)
	env := app.Env
	port := env.ServerPort
	if port == "" {
		port = defaultPort
	}
	// initiate database connection
	database.Init(env)
	redis.NewsCacheService{}.Setup(env)
	//srv := handler.NewDefaultServer()
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", middlewares.EnsureValidToken(env)(srv))
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
	r := gin.Default()
	r.Use(middlewares.GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler(env))
	r.GET("/", playgroundHandler())
	r.Run(":" + port)
}
