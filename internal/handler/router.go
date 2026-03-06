package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "test-junior-go/docs"
)

func NewRouter(h *SubscriptionHandler) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	{
		api.POST("/subscriptions", h.Create)
		api.GET("/subscriptions/:id", h.GetByID)
		api.GET("/subscriptions", h.GetAll)
		api.DELETE("/subscriptions/:id", h.Delete)

		api.GET("/subscriptions/total", h.GetTotal)
	}

	return r
}
