package main

import (
	_ "github.com/deka-microservices/go-user-service/internal/config"
	"github.com/deka-microservices/go-user-service/internal/consts"
	"github.com/deka-microservices/go-user-service/internal/database"
	"github.com/deka-microservices/go-user-service/internal/database/query"
	"github.com/deka-microservices/go-user-service/internal/jwt"
	"github.com/deka-microservices/go-user-service/internal/routes"
	"github.com/rs/zerolog/log"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Version = "dev"

func main() {
	e := gin.New()
	e.Use(logger.SetLogger())
	e.Use(gin.Recovery())

	log.Info().Str("version", Version).Msg("current version")

	// Prepare JWT
	jwtMiddleware, err := jwt.GetJWTMiddleware()
	if err != nil {
		log.Fatal().Err(err).Msg("jwt init")
	}

	// Init DB
	db := database.NewUserDB()
	database.Migrate(db)
	query.SetDefault(db)

	// Init routes
	e.POST("/login", jwtMiddleware.LoginHandler)
	e.POST("/register", routes.Register)
	auth := e.Group("/")
	{
		auth.Use(jwtMiddleware.MiddlewareFunc())
		auth.GET("/refresh_token", jwtMiddleware.RefreshHandler)
	}

	address := viper.GetString(consts.CONFIG_IP) + ":" + viper.GetString(consts.CONFIG_PORT)
	e.Run(address)
}
