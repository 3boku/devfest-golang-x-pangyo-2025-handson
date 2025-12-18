package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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

	// Gemini API 키 확인
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, GenerateBackgroundResponse{
			Success: false,
			Error:   "GEMINI_API_KEY가 설정되지 않았습니다",
		})
		return
	}

	// Gemini 클라이언트 생성
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenerateBackgroundResponse{
			Success: false,
			Error:   "Gemini 클라이언트 생성 실패: " + err.Error(),
		})
		return
	}
	defer client.Close()

	// Gemini 2.0 Flash 모델 사용 (이미지 생성 지원)
	model := client.GenerativeModel("gemini-2.0-flash-exp")

	// 이미지 생성을 위한 설정
	model.ResponseMIMEType = "text/plain"

	// 크리스마스 테마를 강조한 프롬프트 생성
	enhancedPrompt := fmt.Sprintf(
		"Generate a beautiful Christmas card background image based on this description: %s. "+
			"Style: festive, warm, holiday themed, high quality, suitable for a greeting card, 4:3 aspect ratio. "+
			"Do not include any text or letters in the image. Only generate the image, no text response.",
		req.Prompt,
	)

	// 이미지 생성 요청
	resp, err := model.GenerateContent(ctx, genai.Text(enhancedPrompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenerateBackgroundResponse{
			Success: false,
			Error:   "이미지 생성 실패: " + err.Error(),
		})
		return
	}

	// 응답에서 이미지 추출
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			if blob, ok := part.(genai.Blob); ok {
				base64Image := base64.StdEncoding.EncodeToString(blob.Data)
				dataUrl := fmt.Sprintf("data:%s;base64,%s", blob.MIMEType, base64Image)

				c.JSON(http.StatusOK, GenerateBackgroundResponse{
					Success:  true,
					ImageUrl: dataUrl,
				})
				return
			}
		}
	}

	c.JSON(http.StatusInternalServerError, GenerateBackgroundResponse{
		Success: false,
		Error:   "생성된 이미지가 없습니다. 텍스트 응답만 받았습니다.",
	})
}
