package api

import (
	_ "github.com/coffee-realist/TaskManager/TaskPublisher/docs"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/service"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/middleware"
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
// @title Task Publisher API
// @version 1.0
// @description API для управления публикацией задач
// @host localhost:8080
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
		tasks.GET("/get/:project", h.GetFinishedTasks)
		tasks.POST("/publish", h.CreateTask)

	}

	logout := router.Group("/logout")
	logout.Use(middleware.AuthMiddleWare())
	{
		logout.POST("", h.Logout)

	}

	return router
}
