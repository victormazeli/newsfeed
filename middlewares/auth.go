package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"newsfeedbackend/config"
	"newsfeedbackend/redis"
	"newsfeedbackend/utils"
)

func Auth(c *gin.Context, env *config.Env) (interface{}, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	if bearerToken == "" {
		er := errors.New("bearer auth can not be empty")
		return nil, er
	}
	// check if token is not invalidated
	result := redis.NewsCacheService{}.GetAppToken(c, "appToken")
	if result == nil {
		e := errors.New("unauthorized")
		return nil, e
	}
	foundToken := utils.IsTokenInSlice(result, bearerToken)

	if foundToken == true {
		sub, err := utils.ValidateToken(bearerToken, env.JwtKey)
		if err != nil {
			return nil, err
		}
		return sub, nil

	} else {
		err := errors.New("invalid token")
		return nil, err
	}

}
