package api

import (
	"TaskBooker/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Login(c *gin.Context) {
	var loginReq dto.LoginReq

	if err := c.BindJSON(&loginReq); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.services.User.Login(loginReq)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Token.Create(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	c.JSON(http.StatusOK, dto.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "пользователь не авторизован")
		return
	}

	err := h.services.Token.DeleteByUserID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "ошибка выхода")
		return
	}

	c.Status(http.StatusOK)
}

func getUserIDFromContext(c *gin.Context) int {
	if userID, exists := c.Get("userID"); exists {
		return userID.(int)
	}
	return 0
}
