package api

import (
	"fmt"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Refresh godoc
// @Summary Обновление токена
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body dto.TokenReq true "Refresh token"
// @Success 200 {object} dto.TokenResp
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router /refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var tokenReq dto.TokenReq
	if err := c.ShouldBindJSON(&tokenReq); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	newAccessToken, newRefreshToken, err := h.services.Token.Refresh(tokenReq.RefreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, fmt.Sprintf("неверный refresh-токен %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
