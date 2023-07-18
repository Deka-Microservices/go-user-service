package servererrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckInternalServerError(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return true
	}

	return false
}
