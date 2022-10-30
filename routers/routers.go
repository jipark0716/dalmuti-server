package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jipark0716/dalmuti/routers/api"
	"github.com/jipark0716/dalmuti/routers/dalmuti"
)

func InitRouter() (router *gin.Engine) {
	router = gin.New()
	router.Use(gin.Logger())
	router.GET("/health-check", api.HealthCheck)
	router.GET("/dalmuti/play/:id", dalmuti.Play)
	return
}
