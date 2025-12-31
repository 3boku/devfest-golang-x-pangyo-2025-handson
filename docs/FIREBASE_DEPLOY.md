# Firebase Hosting + Cloud Run 배포 가이드

이 문서는 Go로 작성된 크리스마스 카드 생성기 어플리케이션을 **Cloud Run**에 배포하고, **Firebase Hosting**을 통해 연결하여 배포하는 방법을 설명합니다.

이 구성을 사용하면 Firebase의 장점(Custom Domain, CDN, 단순한 배포 명령)과 Cloud Run의 장점(Serverless Container)을 모두 누릴 수 있습니다.

## 1. 사전 준비

### 필수 도구 설치
*   [Google Cloud CLI (gcloud)](https://cloud.google.com/sdk/docs/install)
*   [Firebase CLI](https://firebase.google.com/docs/cli)

### 프로젝트 로그인 및 설정
1.  **GCloud 로그인**: `gcloud auth login`
2.  **Firebase 로그인**: `firebase login`
3.  **프로젝트 설정**:
    ```bash
    # 사용할 Google Cloud 프로젝트 ID 설정
    gcloud config set project [YOUR_PROJECT_ID]
    ```

## 2. Cloud Run 배포

먼저 Go 애플리케이션을 컨테이너로 빌드하고 Cloud Run에 배포합니다.

```bash
# 1. Cloud Run 배포 (소스 코드에서 바로 배포)
gcloud run deploy christmas-card-generator \
  --source . \
  --region asia-northeast3 \
  --allow-unauthenticated \
  --set-env-vars GEMINI_API_KEY=[YOUR_GEMINI_KEY],FIREBASE_PROJECT_ID=[YOUR_PROJECT_ID]
```

*   `--region asia-northeast3`: 서울 리전 사용.
*   `--allow-unauthenticated`: 외부 접속 허용.
*   `--set-env-vars`: 필요한 환경 변수 설정. (`.env` 파일 내용은 배포되지 않으므로 여기서 설정해야 합니다.)
    *   **주의**: `serviceAccountKey.json` 파일은 보안상 이미지에 포함하지 않는 것이 좋습니다. 대신 Cloud Run의 서비스 계정에 적절한 IAM 권한(Firebase Admin SDK 등)을 부여하는 것이 표준입니다.

배포가 완료되면 `Service URL`이 출력됩니다. (예: `https://christmas-card-generator-xxxxx-du.a.run.app`)

## 3. Firebase Hosting 설정

Firebase Hosting을 통해 위에서 배포한 Cloud Run 서비스로 트래픽을 전달하도록 설정합니다.

### 3.1 Firebase 초기화
프로젝트 루트에서 실행합니다.

```bash
firebase init hosting
```
1.  **Project Setup**: 기존 프로젝트 선택.
2.  **Hosting Setup**:
    *   `What do you want to use as your public directory?`: **`static`** (정적 파일이 여기 있다면) 또는 기본값 `public` 사용하고 나중에 설정 수정. 일단 엔터.
    *   `Configure as a single-page app?`: **No**
    *   `Set up automatic builds and deploys with GitHub?`: **No** (나중에 설정 가능)

### 3.2 firebase.json 수정

생성된 `firebase.json` 파일을 열어 아래와 같이 수정합니다. 모든 요청을 Cloud Run 서비스로 다시 씁니다.

```json
{
  "hosting": {
    "public": "static",
    "rewrites": [
      {
        "source": "**",
        "run": {
          "serviceId": "christmas-card-generator",
          "region": "asia-northeast3"
        }
      }
    ],
    "ignore": [
      "firebase.json",
      "**/.*",
      "**/node_modules/**"
    ]
  }
}
```
*   `"serviceId"`: 2번 단계에서 배포한 Cloud Run 서비스 이름.
*   `"region"`: Cloud Run 서비스를 배포한 리전.

## 4. 최종 배포

Firebase Hosting을 배포합니다.

```bash
firebase deploy --only hosting
```

배포가 완료되면 `Hosting URL` (예: `https://[PROJECT_ID].web.app`)을 통해 접속할 수 있습니다.

## 5. (참고) 권한 문제 해결

Cloud Run 앱이 Firebase 서비스(Firestore, Storage)에 접근하려면 권한이 필요합니다.

1.  Cloud Run이 사용하는 기본 서비스 계정 확인:
    `[PROJECT_NUMBER]-compute@developer.gserviceaccount.com`
2.  해당 계정에 IAM 권한 부여:
    *   `Cloud Datastore User` (Firestore 접근)
    *   `Storage Object Admin` (Storage 접근)

GCP 콘솔의 [IAM 및 관리] 에서 권한을 추가할 수 있습니다.
