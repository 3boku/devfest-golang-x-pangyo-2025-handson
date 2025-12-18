# Go와 나노바나나로 크리스마스 카드 작성 Web Application
Go의 강력한 성능과 나노바나나를 활용해, 세상에 하나뿐인 크리스마스 카드 생성 웹 앱을 직접 구현해봅니다.
서버 구축부터 카드 생성 로직까지 단계별로 진행되는 핸즈온을 통해 Go와 Gemini를 즐겨보세요.
이번 크리스마스에는 내 코드로 만든 특별한 카드로 소중한 사람들에게 따뜻한 마음을 전할 수 있습니다.

# ⚙️ 개발 환경
- `Go 1.25.2`
- `HTML5 / CSS3 / Javascript`
- **IDE** : Antigravity
- **Framework**
    - github.com/gin-gonic/gin v1.9.1
    - github.com/google/generative-ai-go v0.19.0

## env
```dotenv
GEMINI_API_KEY="Gemini API Key"
```
```shell
cp .env.example .env
```
`.env.exapmle`파일을 카피하여 `.env` 파일을 만들어 주세요.

# 사용방법
```shell
go get
go run main.go
```

# System Instructions
```
Generate a beautiful Christmas card background image based on this description
Style: festive, warm, holiday themed, high quality, suitable for a greeting card, 4:3 aspect ratio.
Do not include any text or letters in the image. Only generate the image, no text response.
```

# 행사 정보
https://www.ticketa.co/events/67

# Presentation
https://docs.google.com/presentation/d/1sNFp9WPLpKFmUjRoquffd1gX3wYLDt89QXHGlWQx59g/edit?usp=sharing
