<div align="center">
<h1>🌐 docs-translater-gemini</h1>

<img src="https://img.shields.io/badge/Go-1.23.3-blue.svg" alt="Go Version">

<p> 마크다운 파일을 HTML로 변환하고 한국어로 번역해주는 도구입니다. Gemini API를 활용하여 번역을 수행합니다.</p>
</div>

## 주요 기능 ✨
* 디렉토리 내 마크다운 파일 재귀적 번역
* 마크다운을 HTML 형식으로 변환
* 다크 테마와 블루버드 컬러 스킴 적용

## 사용 방법 📝

1. Gemini API 토큰을 `GEMINI_TOKEN` 환경 변수로 설정하세요.
2. 다음과 같이 실행하세요:

```bash
./docs-translater-gemini <path> <exclude_folder_name>
```
