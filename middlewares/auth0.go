package middlewares

//
//import (
//	"context"
//	"errors"
//	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
//	"github.com/gin-gonic/gin"
//	"log"
//	"net/http"
//	"newsfeedbackend/config"
//	"time"
//
//	"github.com/auth0/go-jwt-middleware/v2/validator"
//)
//
//// CustomClaims contains custom data we want from the token.
//type CustomClaims struct {
//	Scope        string `json:"scope"`
//	Sub          string `json:"sub"`
//	ShouldReject bool   `json:"shouldReject,omitempty"`
//}
//
//func (c CustomClaims) Validate(ctx context.Context) error {
//	if c.ShouldReject {
//		return errors.New("should reject was set to true")
//	}
//	return nil
//}
//
//// CheckJWT checkJWT is a gin.HandlerFunc middleware
//// that will check the validity of our JWT.
//func CheckJWT(env *config.Env) gin.HandlerFunc {
//	var claims CustomClaims
//	// The signing key for the token.
//	signingKey := []byte(env.JwtKey)
//
//	// The issuer of our token.
//	issuer := env.AuthDomain
//
//	// The audience of our token.
//	audience := []string{env.AuthAudience}
//
//	// Our token must be signed using this data.
//	keyFunc := func(ctx context.Context) (interface{}, error) {
//		return signingKey, nil
//	}
//
//	// We want this struct to be filled in with
//	// our custom claims from the token.
//	customClaimsData := func() validator.CustomClaims {
//		return &claims
//	}
//	// Set up the validator.
//	jwtValidator, err := validator.New(
//		keyFunc,
//		validator.HS256,
//		issuer,
//		audience,
//		validator.WithCustomClaims(customClaimsData),
//		validator.WithAllowedClockSkew(30*time.Second),
//	)
//	if err != nil {
//		log.Fatalf("failed to set up the validator: %v", err)
//	}
//
//	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
//		log.Printf("Encountered error while validating token: %v", err)
//	}
//
//	middleware := jwtmiddleware.New(
//		jwtValidator.ValidateToken,
//		jwtmiddleware.WithErrorHandler(errorHandler),
//	)
//
//	return func(ctx *gin.Context) {
//		encounteredError := true
//		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
//			encounteredError = false
//			ctx.Request = r
//			ctx.Next()
//		}
//
//		middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)
//
//		if encounteredError {
//			ctx.AbortWithStatusJSON(
//				http.StatusUnauthorized,
//				map[string]string{"message": "Invalid Token."},
//			)
//		}
//	}
//}
