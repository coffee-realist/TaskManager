package api

import (
	"TaskBooker/internal/domain/dto"
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
	taskResp, err := h.services.Task.GetAllByProject(taskReq, userID)
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

func (h *Handler) BookTask(c *gin.Context) {
	var task dto.TaskResp
	if err := c.BindJSON(&task); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	err := h.services.Task.Book(task, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) FinishTask(c *gin.Context) {
	var taskFinish dto.TaskFinish
	if err := c.BindJSON(&taskFinish); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Task.Finish(taskFinish.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.Status(http.StatusOK)
}
