# 🧑‍💻 코드 설명 및 독학 가이드 (Self-Study Guide)

이 문서는 프로젝트의 주요 기능이 어떻게 구현되어 있는지 설명합니다. 코드를 직접 분석하거나 기능을 확장해볼 때 참고하세요.

## 🎯 왜 Go와 Gemini Nano 인가요?

이 프로젝트는 두 가지 핵심 기술을 경험하기 위해 설계되었습니다.

1.  **Golang (Go)의 심플함과 성능**:
    *   Go는 구문이 간결하여 배우기 쉽고, 컴파일 언어로서 강력한 성능을 제공합니다.
    *   본 프로젝트에서 Go는 정적 파일 서빙, API 라우팅, 그리고 **동시성(Concurrency)**을 이용한 비동기 작업(Firebase 업로드 등)을 효율적으로 처리합니다.

2.  **Gemini Nano (Banana Model)**:
    *   Google의 모바일 및 웹 친화적인 최신 경량화 모델입니다.
    *   빠른 응답 속도와 효율성을 자랑하며, 실시간에 가까운 이미지 생성 경험을 제공합니다.

---

## 📂 프로젝트 구조

```text
.
├── handlers/          # API 요청을 처리하는 핵심 로직
│   ├── gemini.go      # Gemini Nano Banana 연동 (배경 생성)
│   ├── card.go        # 이미지 처리 및 텍스트 합성 (fogleman/gg)
│   ├── firebase.go    # Firebase Storage/Firestore 연동
│   └── gallery.go     # 갤러리 조회 로직
├── templates/         # HTML 템플릿 (Frontend)
├── static/            # CSS, JS, Fonts 등 정적 파일
├── main.go            # 서버 진입점 및 라우팅 설정
└── .env               # 환경 변수 (API Key 등)
```

## 🗝️ 주요 함수 설명

### 1. AI 배경 생성 (`handlers/gemini.go`)

*   **`GenerateBackground` 함수**:
    *   **Gemini Client**: `gemini-3-pro-image-preview` 모델을 호출합니다.
    *   이 모델은 "Banana"라는 코드네임으로 불리며, 빠른 이미지 생성 추론에 최적화되어 있습니다.
    *   **Prompt Engineering**: 사용자의 입력 뒤에 "Style: festive, warm, holiday themed..." 같은 스타일 가이드를 덧붙여 크리스마스 분위기의 이미지가 나오도록 유도합니다.

### 2. 카드 이미지 합성 (`handlers/card.go`)

*   **`CreateCard` 함수**:
    *   생성된 배경 이미지 위에 사용자 메시지를 합성합니다.
*   **`overlayText` 함수**:
    *   `fogleman/gg` 라이브러리를 사용합니다. 이는 Go에서 2D 그래픽을 그리기 위한 강력한 도구입니다. Python의 Pillow와 비슷하지만, Go 특유의 정적 타입 안정성을 제공합니다.
    *   **투명 그래디언트**: 텍스트 가독성을 높이기 위해 하단에 검은색 반투명 그라데이션을 그리는 로직을 확인해보세요.

### 3. Firebase 연동 (`handlers/firebase.go`)

*   **goroutine 활용**:
    *   `card.go`에서 `UploadImageToStorage`를 호출할 때 `go func() { ... }()`를 사용하여 **별도의 고루틴**에서 실행합니다.
    *   덕분에 사용자는 이미지를 기다릴 필요 없이 즉시 응답을 받을 수 있습니다. 이것이 Go의 강력한 동시성 모델입니다.

## 🤓 도전 과제 (Self-Study Challenges)

이 코드를 기반으로 더 발전시켜보고 싶다면 다음 과제에 도전해보세요!

1.  **다양한 폰트 지원하기**:
    *   `fonts/` 디렉토리에 새로운 폰트 파일(.ttf)을 추가하고, 사용자가 프론트엔드에서 폰트를 선택할 수 있도록 `handlers/card.go`를 수정해보세요.
2.  **SNS 공유 기능**:
    *   생성된 카드의 URL을 복사하거나, Twitter/Facebook 공유 버튼을 `templates/index.html`에 추가해보세요. (Meta Tag Open Graph 설정 필요)
3.  **좋아요 기능 추가**:
    *   갤러리에서 마음에 드는 카드에 '좋아요'를 누를 수 있도록 Firestore에 `likes` 필드를 추가하고 API를 구현해보세요.

## 🔗 참고 자료

*   [Google Generative AI Go SDK](https://github.com/google/generative-ai-go)
*   [Go Graphics (gg)](https://github.com/fogleman/gg)
*   [Gin Web Framework](https://github.com/gin-gonic/gin)
