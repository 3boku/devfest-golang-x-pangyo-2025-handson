package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	firebaseApp     *firebase.App
	firestoreClient *firestore.Client
	storageBucket   string
)

// CardMetadata 카드 메타데이터 구조체
type CardMetadata struct {
	ID            string    `firestore:"id" json:"id"`
	ImageUrl      string    `firestore:"imageUrl" json:"imageUrl"`
	Message       string    `firestore:"message" json:"message"`
	CreatedAt     time.Time `firestore:"createdAt" json:"createdAt"`
	BackgroundUrl string    `firestore:"backgroundUrl" json:"backgroundUrl"` // Optional: if needed
}

// InitFirebase Firebase 초기화
func InitFirebase() {
	ctx := context.Background()
	var opts []option.ClientOption

	// 서비스 계정 키 파일이 있으면 사용
	credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credsFile))
	}

	conf := &firebase.Config{
		ProjectID:     os.Getenv("FIREBASE_PROJECT_ID"),
		StorageBucket: os.Getenv("FIREBASE_STORAGE_BUCKET"),
	}

	// If FIREBASE_STORAGE_BUCKET is not set, try to derive it or let user know
	// Usually it's PROJECT_ID.appspot.com or similar.
	// For now, we will rely on env var. If not present, storage features might fail or we can default.
	if conf.StorageBucket == "" && conf.ProjectID != "" {
		conf.StorageBucket = conf.ProjectID + ".appspot.com"
	}
	storageBucket = conf.StorageBucket

	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		log.Printf("Firebase 초기화 경고: %v (갤러리 기능이 제한될 수 있습니다)", err)
		return
	}
	firebaseApp = app

	// Firestore 클라이언트 초기화
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("Firestore 초기화 경고: %v", err)
	} else {
		firestoreClient = client
		log.Println("Firebase Firestore 연결 성공")
	}
}

// UploadImageToStorage 이미지를 Firebase Storage에 업로드
func UploadImageToStorage(ctx context.Context, imageData []byte, filename string) (string, error) {
	if firebaseApp == nil {
		return "", fmt.Errorf("Firebase가 초기화되지 않았습니다")
	}

	client, err := firebaseApp.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("Storage 클라이언트 생성 실패: %v", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("버킷 핸들 가져오기 실패: %v", err)
	}

	obj := bucket.Object(filename)
	wc := obj.NewWriter(ctx)
	wc.ContentType = "image/png"

	// Make public?
	// Easier for this hands-on: Make the object public or use signed URL.
	// We will try running specific ACL logic later if needed, but for now standard upload.
	// Note: Default buckets usually enforce strict ACLs.
	// Let's create a Signed URL or just return the standard public URL if the bucket is configured public.
	// For simplicity, we will attempt to make this specific object public if we can,
	// OR we assume the bucket allows public read.
	// However, without making it public, we can't access it via simple URL unless we generate a signed one.

	if _, err := wc.Write(imageData); err != nil {
		return "", fmt.Errorf("이미지 쓰기 실패: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer 닫기 실패: %v", err)
	}

	// Make public explicitely (requires permission)
	// if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
	// 	 log.Println("Warning: Failed to make object public:", err)
	// }

	// Construct public URL
	// Format: https://storage.googleapis.com/<bucket>/<object>
	// Or Firebase Storage download URL format.
	// For simplicity let's return the googleapis URL.
	publicUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s", storageBucket, filename)

	return publicUrl, nil
}

// SaveCardMetadata 카타 메타데이터를 Firestore에 저장
func SaveCardMetadata(ctx context.Context, card CardMetadata) error {
	if firestoreClient == nil {
		return fmt.Errorf("Firestore가 연결되지 않았습니다")
	}

	_, err := firestoreClient.Collection("cards").Doc(card.ID).Set(ctx, card)
	return err
}

// GetRecentCards 최근 카드 목록 조회
func GetRecentCards(ctx context.Context, limit int) ([]CardMetadata, error) {
	if firestoreClient == nil {
		return nil, fmt.Errorf("Firestore가 연결되지 않았습니다")
	}

	iter := firestoreClient.Collection("cards").
		OrderBy("createdAt", firestore.Desc).
		Limit(limit).
		Documents(ctx)

	var cards []CardMetadata
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var card CardMetadata
		if err := doc.DataTo(&card); err != nil {
			continue
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// ParseMultipartFile multipart/form-data에서 파일 추출 (Gallery upload support if needed, mostly internal utils)
func ParseMultipartFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
