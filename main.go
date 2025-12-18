package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 정적 파일 서빙
	r.Static("/static", "./static")

	// HTML 템플릿 로드
	r.LoadHTMLGlob("templates/*")

	// 메인 페이지
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 서버 시작
	log.Println("🎄 크리스마스 카드 생성기 서버 시작: http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("서버 시작 실패:", err)
	}
}
