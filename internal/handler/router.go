package handler

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(h *SubscriptionHandler) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())

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