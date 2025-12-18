// 현재 배경 이미지 URL 저장
let currentBackgroundUrl = '';

// 로딩 표시
function showLoading() {
    document.getElementById('loading').style.display = 'flex';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

// 배경 이미지 생성
async function generateBackground() {
    const prompt = document.getElementById('imagePrompt').value.trim();

    if (!prompt) {
        alert('카드 배경 설명을 입력해주세요!');
        return;
    }

    showLoading();

    try {
        const response = await fetch('/api/generate-background', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ prompt: prompt })
        });

        const data = await response.json();

        if (data.success) {
            currentBackgroundUrl = data.imageUrl;
            const preview = document.getElementById('cardPreview');
            preview.innerHTML = `<img src="${data.imageUrl}" alt="카드 배경">`;
            document.getElementById('createCardBtn').disabled = false;
        } else {
            alert('배경 생성 실패: ' + data.error);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('오류가 발생했습니다: ' + error.message);
    } finally {
        hideLoading();
    }
}

// 카드 완성하기
async function createCard() {
    const message = document.getElementById('cardMessage').value.trim();

    if (!message) {
        alert('카드 메시지를 입력해주세요!');
        return;
    }

    if (!currentBackgroundUrl) {
        alert('먼저 배경을 생성해주세요!');
        return;
    }

    showLoading();

    try {
        const response = await fetch('/api/create-card', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                backgroundUrl: currentBackgroundUrl,
                message: message
            })
        });

        const data = await response.json();

        if (data.success) {
            document.getElementById('finalCard').src = data.cardUrl;
            document.getElementById('resultSection').style.display = 'block';

            // 결과 섹션으로 스크롤
            document.getElementById('resultSection').scrollIntoView({ behavior: 'smooth' });
        } else {
            alert('카드 생성 실패: ' + data.error);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('오류가 발생했습니다: ' + error.message);
    } finally {
        hideLoading();
    }
}

// 카드 다운로드
function downloadCard() {
    const img = document.getElementById('finalCard');
    const link = document.createElement('a');
    link.href = img.src;
    link.download = 'christmas-card.png';
    link.click();
}
