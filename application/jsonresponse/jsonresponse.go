package jsonresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"result":  data,
	})
}

func Failed(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "failed",
		"result":  data,
	})
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "error",
		"result":  err,
	})
}
