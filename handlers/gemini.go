package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GenerateBackgroundRequest 배경 생성 요청
type GenerateBackgroundRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// GenerateBackgroundResponse 배경 생성 응답
type GenerateBackgroundResponse struct {
	Success  bool   `json:"success"`
	ImageUrl string `json:"imageUrl,omitempty"`
	Error    string `json:"error,omitempty"`
}

// GenerateBackground Gemini를 사용해 배경 이미지 생성
func GenerateBackground(c *gin.Context) {
	var req GenerateBackgroundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenerateBackgroundResponse{
			Success: false,
			Error:   "잘못된 요청입니다: " + err.Error(),
		})
		return
	}

	// gemini 호출
}
