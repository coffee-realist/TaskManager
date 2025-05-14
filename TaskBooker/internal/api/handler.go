package api

import (
	_ "github.com/coffee-realist/TaskManager/TaskBooker/docs"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/domain/service"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// InitRoutes godoc
// @title Task Booker API
// @version 1.0
// @description API для управления задачами в микросервисе Booker
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	unauthorized := router.Group("/")
	{
		unauthorized.POST("/login", h.Login)
		unauthorized.POST("/refresh", h.Refresh)
	}

	tasks := router.Group("/tasks")
	tasks.Use(middleware.AuthMiddleWare())
	{
		tasks.GET("/get/:project", h.GetTasks)
		tasks.POST("/book", h.BookTask)
		tasks.POST("/finish", h.FinishTask)

	}
	logout := router.Group("/logout")
	logout.Use(middleware.AuthMiddleWare())
	{
		logout.POST("", h.Logout)

	}
	return router
}
