package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetGallery 최근 공유된 카드 목록을 가져와서 렌더링
func GetGallery(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Firebase에서 최근 카드 가져오기
	cards, err := GetRecentCards(ctx, 20) // 최근 20개
	if err != nil {
		// 에러 발생 시 (설정 안됨 등) 빈 목록으로 렌더링하거나 에러 페이지
		// 여기서는 빈 목록으로 처리하고 템플릿에서 안내
		cards = []CardMetadata{}
	}

	c.HTML(http.StatusOK, "gallery.html", gin.H{
		"Cards": cards,
	})
}
