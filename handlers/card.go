package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // JPEG 디코딩 지원
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"golang.org/x/image/font/basicfont"
)

// CreateCardRequest 카드 생성 요청
type CreateCardRequest struct {
	BackgroundUrl string `json:"backgroundUrl" binding:"required"`
	Message       string `json:"message" binding:"required"`
}

// CreateCardResponse 카드 생성 응답
type CreateCardResponse struct {
	Success bool   `json:"success"`
	CardUrl string `json:"cardUrl,omitempty"`
	Error   string `json:"error,omitempty"`
}

// CreateCard 배경 이미지에 메시지를 추가하여 카드 생성
func CreateCard(c *gin.Context) {
	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CreateCardResponse{
			Success: false,
			Error:   "잘못된 요청입니다: " + err.Error(),
		})
		return
	}

	// Base64 이미지 디코딩
	bgImage, err := decodeBase64Image(req.BackgroundUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreateCardResponse{
			Success: false,
			Error:   "이미지 디코딩 실패: " + err.Error(),
		})
		return
	}

	// 이미지에 텍스트 오버레이
	cardImage, err := overlayText(bgImage, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreateCardResponse{
			Success: false,
			Error:   "텍스트 오버레이 실패: " + err.Error(),
		})
		return
	}

	// 결과 이미지를 Base64로 인코딩
	var buf bytes.Buffer
	if err := png.Encode(&buf, cardImage); err != nil {
		c.JSON(http.StatusInternalServerError, CreateCardResponse{
			Success: false,
			Error:   "이미지 인코딩 실패: " + err.Error(),
		})
		return
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", base64Image)

	c.JSON(http.StatusOK, CreateCardResponse{
		Success: true,
		CardUrl: dataUrl,
	})
}

// decodeBase64Image Base64 데이터 URL에서 이미지 디코딩
func decodeBase64Image(dataUrl string) (image.Image, error) {
	// data:image/png;base64,... 형식에서 Base64 부분 추출
	parts := strings.SplitN(dataUrl, ",", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("잘못된 데이터 URL 형식")
	}

	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Base64 디코딩 실패: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("이미지 디코딩 실패: %w", err)
	}

	return img, nil
}

// overlayText 이미지 위에 텍스트를 오버레이
func overlayText(bgImage image.Image, message string) (image.Image, error) {
	bounds := bgImage.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 새로운 컨텍스트 생성
	dc := gg.NewContext(width, height)

	// 배경 이미지 그리기
	dc.DrawImage(bgImage, 0, 0)

	// 폰트 크기 계산 (이미지 크기에 비례)
	fontSize := float64(height) / 15
	if fontSize < 24 {
		fontSize = 24
	}
	if fontSize > 72 {
		fontSize = 72
	}

	// TrueType 폰트 로드 시도
	fontLoaded := false
	fontPath := "./fonts/NotoSansKR-Bold.ttf"
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Printf("폰트 로드 실패 (%s): %v, 기본 폰트 사용", fontPath, err)
		// 기본 폰트 사용
		dc.SetFontFace(basicfont.Face7x13)
		fontSize = 13
	} else {
		fontLoaded = true
	}

	// 텍스트 영역 (하단에 반투명 배경)
	textBoxHeight := float64(height) * 0.28
	textBoxY := float64(height) - textBoxHeight

	// 반투명 배경 그라데이션 효과
	for i := 0; i < int(textBoxHeight); i++ {
		alpha := 0.7 * float64(i) / textBoxHeight
		dc.SetRGBA(0, 0, 0, alpha)
		y := textBoxY + float64(i)
		dc.DrawLine(0, y, float64(width), y)
		dc.Stroke()
	}

	// 폰트가 로드된 경우 다시 설정 (그라데이션 그린 후)
	if fontLoaded {
		dc.LoadFontFace(fontPath, fontSize)
	} else {
		dc.SetFontFace(basicfont.Face7x13)
	}

	// 메시지를 줄바꿈 처리
	maxWidth := float64(width) * 0.85
	lines := wrapText(dc, message, maxWidth)

	// 텍스트 설정 - 그림자 효과
	lineHeight := fontSize * 1.4
	if !fontLoaded {
		lineHeight = 20 // 기본 폰트용 줄 높이
	}

	// 텍스트 시작 Y 위치 계산 (수직 중앙 정렬)
	totalTextHeight := float64(len(lines)) * lineHeight
	startY := textBoxY + (textBoxHeight-totalTextHeight)/2 + fontSize

	// 그림자 그리기
	dc.SetRGBA(0, 0, 0, 0.5)
	for i, line := range lines {
		y := startY + float64(i)*lineHeight
		dc.DrawStringAnchored(line, float64(width)/2+2, y+2, 0.5, 0.5)
	}

	// 메인 텍스트 그리기 (흰색)
	dc.SetRGB(1, 1, 1)
	for i, line := range lines {
		y := startY + float64(i)*lineHeight
		dc.DrawStringAnchored(line, float64(width)/2, y, 0.5, 0.5)
	}

	return dc.Image(), nil
}

// wrapText 텍스트를 주어진 너비에 맞게 줄바꿈
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string

	// 먼저 줄바꿈 문자로 분리
	paragraphs := strings.Split(text, "\n")

	for _, para := range paragraphs {
		words := strings.Fields(para)

		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		currentLine := words[0]

		for _, word := range words[1:] {
			testLine := currentLine + " " + word
			w, _ := dc.MeasureString(testLine)

			if w <= maxWidth {
				currentLine = testLine
			} else {
				lines = append(lines, currentLine)
				currentLine = word
			}
		}

		lines = append(lines, currentLine)
	}

	return lines
}
