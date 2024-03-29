package main

import (
	"context"
	"flag"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"newsfeedbackend/config"
	"newsfeedbackend/database"
	"newsfeedbackend/graph"
	"newsfeedbackend/graph/generated"
	"newsfeedbackend/middlewares"
	"newsfeedbackend/redis"
	"newsfeedbackend/utils"
)

const defaultPort = "8080"

// Defining the Graphql handler
func graphqlHandler(env *config.Env) gin.HandlerFunc {
	// pass the env to the resolver since the resolver instance is where we can initialize dependencies
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			Env: env,
		},
	}))
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
	// Toggle between different environments using command line
	environment := flag.String("e", "development", "")
	flag.Parse()

	// Initialize the application
	app := config.App(*environment)
	env := app.Env
	port := env.ServerPort
	if port == "" {
		port = defaultPort
	}

	// Initialize database connection
	database.Init(env)

	// Initialize redis service
	redis.NewsCacheService{}.Setup(env)

	// Initialize cron job
	utils.InitCron(context.Background(), env)

	// Create the Gin router
	if *environment == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler(env))
	r.GET("/", playgroundHandler())

	// Start the server
	r.Run(":" + port)
}
