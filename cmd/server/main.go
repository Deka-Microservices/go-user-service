package main

import (
	_ "github.com/deka-microservices/go-user-service/internal/config"
	"github.com/deka-microservices/go-user-service/internal/consts"
	"github.com/deka-microservices/go-user-service/internal/database"
	"github.com/deka-microservices/go-user-service/internal/database/query"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	e := gin.New()
	e.Use(logger.SetLogger())
	e.Use(gin.Recovery())

	// Init DB
	db := database.NewUserDB()
	database.Migrate(db)
	query.SetDefault(db)

	address := viper.GetString(consts.CONFIG_IP) + ":" + viper.GetString(consts.CONFIG_PORT)
	e.Run(address)
}
