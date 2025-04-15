package api

import (
	"TaskPublisher/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (h *Handler) GetTasks(c *gin.Context) {
	var taskReq dto.TaskReq
	if err := c.BindJSON(&taskReq); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	taskResp, err := h.services.Task.GetAllFinishedByProject(taskReq, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		select {
		case task, ok := <-taskResp:
			if !ok {
				return false
			}
			c.SSEvent("message", task)
			return true
		case <-c.Writer.CloseNotify():
			return false
		}
	})
}

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
