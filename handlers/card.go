package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
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

	// 텍스트 영역 (하단에 반투명 배경)
	textBoxHeight := float64(height) * 0.25
	textBoxY := float64(height) - textBoxHeight

	// 반투명 검은색 배경
	dc.SetRGBA(0, 0, 0, 0.6)
	dc.DrawRectangle(0, textBoxY, float64(width), textBoxHeight)
	dc.Fill()

	// 기본 폰트 사용 (gg 라이브러리의 기본 폰트 대신 basicfont 사용)
	face := basicfont.Face7x13
	dc.SetFontFace(face)

	// 텍스트 설정
	dc.SetRGB(1, 1, 1) // 흰색

	// 텍스트 크기 계산을 위한 스케일 (기본 폰트가 작으므로 여러 줄로 표시)
	fontSize := 13.0 // basicfont.Face7x13의 높이
	lineHeight := fontSize * 1.5
	maxWidth := float64(width) * 0.9
	padding := float64(width) * 0.05

	// 메시지를 줄바꿈 처리
	lines := wrapText(dc, message, maxWidth)

	// 텍스트 시작 Y 위치 계산 (수직 중앙 정렬)
	totalTextHeight := float64(len(lines)) * lineHeight
	startY := textBoxY + (textBoxHeight-totalTextHeight)/2 + fontSize

	// 각 줄 그리기
	for i, line := range lines {
		y := startY + float64(i)*lineHeight
		// 수평 중앙 정렬
		w, _ := dc.MeasureString(line)
		x := (float64(width) - w) / 2
		if x < padding {
			x = padding
		}
		dc.DrawString(line, x, y)
	}

	return dc.Image(), nil
}

// wrapText 텍스트를 주어진 너비에 맞게 줄바꿈
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string
	words := strings.Fields(text)

	if len(words) == 0 {
		return lines
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
	return lines
}
