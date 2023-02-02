package config

import (
	"github.com/spf13/viper"
	"log"
)

// Env We declare an env struct which will serve as model to map our env to
type Env struct {
	NewsApiBaseUrl string `mapstructure:"NEWS_API_BASE_URL"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPass         string `mapstructure:"DB_PASS"`
	DBName         string `mapstructure:"DB_NAME"`
}

// NewEnv we use viper as our configuration tool to read and load our configurations from env file
func NewEnv(envType string) *Env {
	env := Env{}
	viper.SetConfigType("env")
	viper.SetConfigName(envType)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file development.env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
