package jwt

import (
	"crypto/rand"
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/deka-microservices/go-user-service/internal/bcrypt"
	"github.com/deka-microservices/go-user-service/internal/consts"
	"github.com/deka-microservices/go-user-service/internal/database/query"
	"github.com/deka-microservices/go-user-service/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const JWT_SECRET_DEFAULT_SIZE int = 512
const IDENTITY_KEY = "username"

func GetJWTSecret() []byte {

	jwt_secret := []byte(viper.GetString(consts.CONFIG_JWT_SECRET))
	// Generate random
	if len(jwt_secret) == 0 {
		jwt_secret = generateJWTKey()
	}

	return []byte(jwt_secret)
}

func generateJWTKey() []byte {
	new_key := make([]byte, JWT_SECRET_DEFAULT_SIZE)

	_, err := rand.Read(new_key)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate jwt key")
	}

	return new_key
}

func GetJWTMiddleware() (*jwt.GinJWTMiddleware, error) {

	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "deka.space",
		Key:             GetJWTSecret(),
		Timeout:         time.Minute,
		MaxRefresh:      time.Minute,
		IdentityKey:     IDENTITY_KEY,
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authenticator:   authenticator,
		Authorizator:    authorization,
		Unauthorized:    unauthorized,
	})

	if err != nil {
		return nil, err
	}

	err = jwtMiddleware.MiddlewareInit()
	return jwtMiddleware, err
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.UserClaims); ok {
		return jwt.MapClaims{
			IDENTITY_KEY: v.Username,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models.UserClaims{
		Username: claims[IDENTITY_KEY].(string),
	}
}

func authenticator(c *gin.Context) (interface{}, error) {
	var login models.UserLoginInfo
	if err := c.ShouldBind(&login); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	user, err := query.User.GetByUsername(login.Username)
	if err != nil {
		// If username not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", jwt.ErrFailedAuthentication
		}
		// otherwise
		log.Error().Err(err).Msg("failed to get user by username")
		return "", err
	}

	password_hash := user.Password

	client, err := bcrypt.NewGrpcBcryptServiceClient()
	if err != nil {
		return nil, err
	}

	if ok, err := client.CheckPassword(c.Request.Context(), login.Password, password_hash); !ok {
		if err != nil {
			log.Error().Err(err).Msg("failed to check password")
		}
		return "", jwt.ErrFailedAuthentication
	}

	return &models.UserClaims{Username: login.Username}, nil
}

func authorization(data interface{}, c *gin.Context) bool {
	if _, ok := data.(*models.UserClaims); ok {
		//TODO! Check rights
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
