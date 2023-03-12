package utils_old

import (
	"context"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"newsfeedbackend/config"
	"newsfeedbackend/graph/model"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrNoToken      = errors.New("no token found")
	ErrTokenExpired = errors.New("token expired")
)

func ValidateToken(bearerToken string, env *config.Env) (*model.User, error) {
	//ar := jwk.NewAutoRefresh(context.Background())
	//ar.Configure(fmt.Sprintf("%s.well-known/jwks.json", env.AuthDomain))
	keySet, er := jwk.Fetch(context.Background(), fmt.Sprintf("%s.well-known/jwks.json", env.AuthDomain))
	if er != nil {
		return nil, er
	}

	token, err := jwt.Parse([]byte(bearerToken), jwt.WithKeySet(keySet))
	if err != nil {
		log.Print("error parsing token")
		return nil, ErrInvalidToken
	}

	if err := jwt.Validate(token,
		jwt.WithAudience(env.AuthAudience),
		jwt.WithIssuer(env.AuthDomain)); err != nil {
		log.Print("error validating token")
		return nil, ErrInvalidToken
	}

	subValue, ok := token.Get("sub")
	if !ok {
		log.Print("error no sub found")
		return nil, ErrInvalidToken
	}

	sub, ok := subValue.(string)
	if !ok {
		log.Print("error sub not a string")
		return nil, ErrInvalidToken
	}

	if token.Expiration().Local().Unix() < time.Now().Local().Unix() {
		log.Print("token expired")
		return nil, ErrTokenExpired
	}

	//payload := make(map[string]interface{})
	//
	//payload["sub"] = struct {
	//	UserId string `json:"user_id"`
	//}{sub}
	log.Print(sub)

	payload := &model.User{
		Email: sub,
	}

	return payload, nil
}
