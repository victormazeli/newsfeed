package middlewares

import (
	"github.com/gin-gonic/gin"
	"newsfeedbackend/config"
	"newsfeedbackend/graph/model"
	"newsfeedbackend/utils"
)

func Auth(c *gin.Context, env *config.Env) (*model.User, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	sub, err := utils.ValidateToken(bearerToken, env)
	if err != nil {
		return nil, err
	}
	return sub, nil

}
