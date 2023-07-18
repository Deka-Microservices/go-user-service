package routes

import (
	"net/http"

	"github.com/deka-microservices/go-user-service/internal/bcrypt"
	"github.com/deka-microservices/go-user-service/internal/database/query"
	servererrors "github.com/deka-microservices/go-user-service/internal/server_errors"
	"github.com/deka-microservices/go-user-service/pkg/models"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user_info models.UserRegisterInfo
	if err := ctx.ShouldBind(&user_info); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	client, err := bcrypt.NewGrpcBcryptServiceClient()
	if servererrors.CheckInternalServerError(ctx, err) {
		return
	}

	hash, err := client.HashPassword(ctx.Request.Context(), user_info.Password)
	if servererrors.CheckInternalServerError(ctx, err) {
		return
	}

	user := models.User{
		Username: user_info.Username,
		Password: hash,
	}

	// TODO! Check if user exists here

	// ------------------------------

	err = query.User.Create(&user)
	if servererrors.CheckInternalServerError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, user)

}
