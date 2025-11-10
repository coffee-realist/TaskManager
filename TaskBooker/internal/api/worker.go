package api

import (
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

// GetTasks godoc
// @Summary Получить задачи по проекту
// @Tags tasks
// @Accept  json
// @Produce  text/event-stream
// @Param project path string true "Название проекта" Example(project1)
// @Success 200 {object} dto.TaskResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security BearerAuth
// @Router /tasks/get/{project} [get]
func (h *Handler) GetTasks(c *gin.Context) {
	var taskReq dto.TaskReq
	if err := c.ShouldBindUri(&taskReq); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	taskResp, err := h.services.Task.GetAllByProject(c.Request.Context(), taskReq)
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

// BookTask godoc
// @Summary Забронировать задачу
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param input body dto.TaskResp true "Данные задачи"
// @Success 200 {object} dto.TaskFinish
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security BearerAuth
// @Router /tasks/book [post]
func (h *Handler) BookTask(c *gin.Context) {
	var task dto.TaskResp
	var taskFinish dto.TaskFinish
	if err := c.BindJSON(&task); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	taskID, err := h.services.Task.Book(task, userID)
	taskFinish.ID = taskID
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, taskFinish)
}

// FinishTask godoc
// @Summary Завершить задачу
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param input body dto.TaskFinish true "ID задачи"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security BearerAuth
// @Router /tasks/finish [post]
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
