package api

import "github.com/gin-gonic/gin"

func HealthCheck(context *gin.Context) {
	context.JSON(200, map[string]interface{}{
		"code": 200,
	})
}
