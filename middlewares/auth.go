package middlewares

import (
	"github.com/gin-gonic/gin"
	"newsfeedbackend/config"
	"newsfeedbackend/utils"
)

func Auth(c *gin.Context, env *config.Env) (interface{}, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	sub, err := utils.ValidateToken(bearerToken, env.JwtKey)
	if err != nil {
		return nil, err
	}
	return sub, nil

}
