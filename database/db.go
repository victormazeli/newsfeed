package database

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"newsfeedbackend/config"
)

func Init(env *config.Env) {
	err := mgm.SetDefaultConfig(nil, env.DBName, options.Client().ApplyURI(env.ConnectionUrl))

	if err != nil {
		panic("Error occurred connecting to database!")
	}
	log.Print("Connected to database!")
}
