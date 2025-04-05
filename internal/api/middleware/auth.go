package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/vcsfrl/xm/dto"
	"github.com/vcsfrl/xm/internal/config"
	"time"
)

const identityKey = "id"

type AuthenticationManager struct {
	AuthMiddleware *jwt.GinJWTMiddleware
	config         *config.Config
	logger         zerolog.Logger
}

func NewAuthenticationManager(config *config.Config, logger zerolog.Logger) (*AuthenticationManager, error) {
	var err error

	result := &AuthenticationManager{config: config, logger: logger}
	result.AuthMiddleware, err = jwt.New(result.buildMiddleware())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (am *AuthenticationManager) JwtHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := am.AuthMiddleware.MiddlewareInit()
		if err != nil {
			am.logger.Error().Err(err).Msg("authMiddleware.MiddlewareInit() Error")
		}
	}
}

func (am *AuthenticationManager) buildMiddleware() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(am.config.AuthJwtSecret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: am.payloadFunc(),

		IdentityHandler: am.identityHandler(),
		Authenticator:   am.authenticator(),
		Authorizator:    am.authorizator(),
		Unauthorized:    am.unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	}
}

func (am *AuthenticationManager) payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*dto.AuthUser); ok {
			return jwt.MapClaims{
				identityKey: v.Username,
			}
		}
		return jwt.MapClaims{}
	}
}

func (am *AuthenticationManager) identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &dto.AuthUser{
			Username: claims[identityKey].(string),
		}
	}
}

// authenticator is the function that checks if the user is authenticated
func (am *AuthenticationManager) authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals dto.LoginRequest
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userName := loginVals.Username
		password := loginVals.Password

		if userName == am.config.AuthUser && password == am.config.AuthPassword {
			return &dto.AuthUser{
				Username: userName,
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

// authorizator is the function that checks if the user is authorized to access the resource
func (am *AuthenticationManager) authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*dto.AuthUser); ok && v.Username == am.config.AuthUser {
			return true
		}
		return false
	}
}

// unauthorized is the function that handles unauthorized access
func (am *AuthenticationManager) unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}
