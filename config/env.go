package config

import (
	"github.com/spf13/viper"
	"log"
)

// Env We declare an env struct which will serve as model to map our env to
type Env struct {
	NewsApiBaseUrl string `mapstructure:"NEWS_API_BASE_URL"`
	NewsApiKey     string `mapstructure:"NEWS_API_KEY"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	ConnectionUrl  string `mapstructure:"CONNECTION_URL"`
	DBName         string `mapstructure:"DB_NAME"`
	AuthDomain     string `mapstructure:"AUTH_DOMAIN"`
	AuthAudience   string `mapstructure:"AUTH_AUDIENCE"`
	JwtKey         string `mapstructure:"JWT_KEY"`
	RedisAddr      string `mapstructure:"REDIS_ADDR"`
	MailGunDomain  string `mapstructure:"MAILGUN_DOMAIN"`
	MailGunApiKey  string `mapstructure:"MAILGUN_API_KEY"`
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
