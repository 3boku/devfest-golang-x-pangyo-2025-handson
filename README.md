# 🎄 DevFest 2025 Pangyo: Golang Hands-on

**DevFest 2025 Pangyo**의 핸즈온 세션, **"나만의 AI 크리스마스 카드 만들기"** 레포지토리입니다.

이 세션의 핵심 목표는 구글의 최신 기술인 **Gemini Nano** (Banana Model)와 강력한 백엔드 언어인 **Golang (Go)** 을 활용하여 실용적인 AI 어플리케이션을 직접 만들어보는 것입니다.

> [!WARNING]
> **비용 발생 주의 (Billing Warning)**
> *   **Gemini API**: 이미지 생성 모델(`gemini-3-pro-image-preview`)은 **Free Tier에서 사용할 수 없습니다.** 반드시 결제 계정이 연결된 프로젝트의 API Key를 사용해야 합니다. (유료 과금 발생)
> *   **Firebase**: Spark 요금제(무료)는 제한된 용량 내에서 무료입니다. 하지만 앱 사용자가 많아져 저장 용량(Storage)이나 읽기/쓰기(Firestore, Hosting) 트래픽이 늘어나면 Blaze 요금제(종량제) 전환 시 비용이 발생할 수 있습니다.
> *   실습 후 불필요한 리소스는 정리(삭제)하는 것을 권장합니다.

## 🎯 세션 목표 (Session Goals)

*   **Golang 경험하기**: 간결하고 강력한 Go 언어를 사용하여 웹 서버를 구축합니다.
*   **Gemini Nano (Banana) 활용**: 구글의 최신 경량화 AI 모델인 **Banana** (Gemini 3.0 Pro Image Preview)를 사용하여 고품질 이미지를 생성합니다.
*   **실제 서비스 구현**: Firebase를 연동하여 결과물을 저장하고 공유하는 전체 서비스 흐름을 이해합니다.

## ✨ 주요 기능

*   **AI 배경 생성 (Gemini Nano Banana)**: "눈 내리는 오두막"과 같은 텍스트만 입력하면 Gemini Banana 모델이 순식간에 아름다운 배경을 그려줍니다.
*   **카드 커스터마이징 (Go Graphics)**: Go의 `fogleman/gg` 라이브러리를 사용하여 텍스트와 이미지를 정교하게 합성합니다.
*   **공유 갤러리**: 내가 만든 카드를 Firebase에 저장하고 친구들과 공유합니다.

> [!NOTE]
> *   **[Firebase 배포 가이드](docs/FIREBASE_DEPLOY.md)**: 만든 앱을 클라우드(Cloud Run)에 배포하고 싶으신가요?
> *   **[코드 설명서](docs/CODE_GUIDE.md)**: Go 코드 구조와 Gemini 연동 방식이 궁금하신가요? (독학용 가이드)

## 🛠️ 기술 스택 (Tech Stack)

*   **Language**: **Go 1.23+** (Simplicity & Performance)
*   **AI Model**: **Google Gemini 3.0 Pro** (Codename: Banana)
    *   *Note: 이 핸즈온에서는 `gemini-3-pro-image-preview` 모델을 사용합니다.*
*   **Image Processing**: `fogleman/gg`
*   **Database & Storage**: Firebase Firestore & Storage

## 🚀 시작하기

### 1. 필수 조건

*   [Go](https://go.dev/dl/) 설치 (1.23 이상 권장)
*   [Google AI Studio](https://aistudio.google.com/)에서 API Key 발급
*   (선택) Firebase 프로젝트 생성 및 `serviceAccountKey.json` 준비

### 2. 설치 및 실행

1.  레포지토리 클론
    ```bash
    git clone https://github.com/curogom/devfest-golang-x-pangyo-2025-handson.git
    cd devfest-golang-x-pangyo-2025-handson
    ```

2.  환경 변수 설정
    `.env.example` 파일을 복사하여 `.env` 파일을 생성하고 키를 입력합니다.
    ```bash
    cp .env.example .env
    ```
    ```env
    # .env
    GEMINI_API_KEY=your_api_key_here
    
    # Firebase 설정 (선택 사항 - 갤러리 기능 사용 시 필요)
    FIREBASE_PROJECT_ID=your_project_id
    GOOGLE_APPLICATION_CREDENTIALS=serviceAccountKey.json
    ```

3.  의존성 설치
    ```bash
    go mod download
    ```

3.5 (선택사항) Firebase 갤러리 활성화
    Firebase 공유 기능을 사용하려면 `main.go` 파일에서 다음 줄의 주석을 해제하세요.
    ```go
    // main.go
    // handlers.InitFirebase()  <-- 주석 해제
    // ...
    // r.GET("/gallery", handlers.GetGallery) <-- 주석 해제
    ```

4.  서버 실행
    ```bash
    go run main.go
    ```

5.  접속
    브라우저에서 [http://localhost:8080](http://localhost:8080) 으로 접속하세요.

## ☁️ 배포 (Deployment)

`Dockerfile`이 포함되어 있어 컨테이너 환경에 쉽게 배포할 수 있습니다.

```bash
docker build -t christmas-card-generator .
docker run -p 8080:8080 --env-file .env christmas-card-generator
```
