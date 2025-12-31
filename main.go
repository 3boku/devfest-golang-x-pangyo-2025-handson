package main

import (
	"log"
	"net/http"

	"devfest-golang-x-pangyo-2025-handson/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("env 불러오기 실패: ", err)
	}

	r := gin.Default()

	// [Optional] Firebase 초기화
	// 갤러리 기능을 사용할 때만 주석을 해제하세요.
	// handlers.InitFirebase()

	// 정적 파일 서빙
	r.Static("/static", "./static")
	r.Static("/output", "./output")

	// HTML 템플릿 로드
	r.LoadHTMLGlob("templates/*")

	// 메인 페이지
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// [Optional] 갤러리 페이지
	// 갤러리 기능을 사용할 때만 주석을 해제하세요.
	// r.GET("/gallery", handlers.GetGallery)

	// API 라우트
	api := r.Group("/api")
	{
		api.POST("/generate-background", handlers.GenerateBackground)
		api.POST("/create-card", handlers.CreateCard)
	}

	// 서버 시작
	log.Println("🎄 크리스마스 카드 생성기 서버 시작: http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("서버 시작 실패:", err)
	}
}
