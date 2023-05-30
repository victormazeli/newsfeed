package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"newsfeedbackend/config"
	"newsfeedbackend/utils"
)

func Auth(c *gin.Context, env *config.Env) (interface{}, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	if bearerToken == "" {
		er := errors.New("bearer auth can not be empty")
		return nil, er
	}
	sub, err := utils.ValidateToken(bearerToken, env.JwtKey)
	if err != nil {
		return nil, err
	}
	return sub, nil

}
