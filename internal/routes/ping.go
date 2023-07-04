package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/deka-microservices/go-bcrypt-service/pkg/service"
	"github.com/deka-microservices/go-user-service/internal/consts"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Ping go bcrypt server
func Ping(ctx *gin.Context) {
	address := viper.GetString(consts.CONFIG_BCRYPT_SERVER_ADDRESS)

	var options []grpc.DialOption
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Error().Err(err).Msg("bcrypt service dial fail")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	defer conn.Close()

	client := service.NewBcryptServiceClient(conn)

	client_ctx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	_, err = client.HashPassword(client_ctx, &service.HashRequest{Password: "test"})
	if err != nil {
		log.Error().Err(err).Msg("bcrypt service fail")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

}
