package api

import (
	"TaskPublisher/internal/domain/service"
	"TaskPublisher/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	unauthorized := router.Group("/")
	{
		unauthorized.POST("/login", h.Login)
		unauthorized.POST("/refresh", h.Refresh)
	}

	profile := router.Group("/tasks")
	profile.Use(middleware.AuthMiddleWare())
	{
		profile.GET("/get", h.GetTasks)
		profile.POST("/book", h.CreateTask)

	}
	logout := router.Group("/logout")
	profile.Use(middleware.AuthMiddleWare())
	{
		logout.POST("", h.Logout)

	}
	return router
}
