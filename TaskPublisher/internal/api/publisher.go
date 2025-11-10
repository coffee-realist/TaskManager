package api

import (
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

// GetFinishedTasks godoc
// @Summary Получить завершенные задачи
// @Tags tasks
// @Accept  json
// @Produce  text/event-stream
// @Param project path string true "Название проекта" Example(project1)
// @Success 200 {object} dto.TaskResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security BearerAuth
// @Router /tasks/get/{project} [get]
func (h *Handler) GetFinishedTasks(c *gin.Context) {
	var taskReq dto.TaskReq
	if err := c.ShouldBindUri(&taskReq); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	taskResp, err := h.services.Task.GetAllFinishedByProject(c.Request.Context(), taskReq)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case task, ok := <-taskResp:
			if !ok {
				return false
			}
			c.SSEvent("message", task)
			return true
		case <-ticker.C:
			c.SSEvent("keepalive", "")
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})

}

// CreateTask godoc
// @Summary Создать новую задачу
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param input body dto.TaskResp true "Данные задачи"
// @Success 201
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security BearerAuth
// @Router /tasks/publish [post]
func (h *Handler) CreateTask(c *gin.Context) {
	var task dto.TaskResp
	if err := c.BindJSON(&task); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	task.PublisherID = getUserIDFromContext(c)
	err := h.services.Task.Publish(task)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.Status(http.StatusCreated)
}
