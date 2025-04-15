package api

import (
	"TaskPublisher/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Refresh(c *gin.Context) {
	var tokenReq dto.TokenReq
	if err := c.ShouldBindJSON(&tokenReq); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	newAccessToken, newRefreshToken, err := h.services.Token.Refresh(tokenReq.RefreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "неверный refresh-токен")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
