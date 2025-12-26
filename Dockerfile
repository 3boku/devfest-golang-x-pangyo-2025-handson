# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 필수 패키지 설치 (필요한 경우)
RUN apk add --no-cache git

# 의존성 파일 복사 및 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run Stage
FROM alpine:latest

WORKDIR /app

# 타임존 설정 (선택 사항)
RUN apk add --no-cache tzdata

# 빌드된 바이너리 복사
COPY --from=builder /app/main .

# 정적 파일 및 템플릿 복사
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/fonts ./fonts

# 포트 노출
EXPOSE 8080

# 실행
CMD ["./main"]
