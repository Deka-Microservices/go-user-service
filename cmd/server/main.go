package main

import (
	_ "github.com/deka-microservices/go-user-service/internal/config"
	"github.com/deka-microservices/go-user-service/internal/consts"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	e := gin.New()
	e.Use(logger.SetLogger())
	e.Use(gin.Recovery())

	address := viper.GetString(consts.CONFIG_IP) + ":" + viper.GetString(consts.CONFIG_PORT)
	e.Run(address)
}
