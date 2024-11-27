package main

import (
	"docs-translater-gemini/internal/gemini"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// 재시도 관련 상수 추가
const (
	retryDelay = 5 * time.Second
)

// 재시도 큐를 위한 구조체 정의
type TranslationJob struct {
	path string
}

var (
	client     *gemini.GeminiClient
	waitGroup  sync.WaitGroup
	retryQueue chan TranslationJob
)

func init() {
	godotenv.Load()

	client = gemini.NewGeminiClient(&gemini.GeminiConfig{
		ApiKey: os.Getenv("GEMINI_TOKEN"),
	})
	retryQueue = make(chan TranslationJob, 100) // 버퍼 사이즈 100의 채널 생성

	// 재시도 워커 시작
	go retryWorker()
}

// 재시도 워커 함수 추가
func retryWorker() {
	for job := range retryQueue {
		time.Sleep(retryDelay)

		waitGroup.Add(1)
		go func(j TranslationJob) {
			defer waitGroup.Done()
			err := translateFile(j.path, true)
			if err != nil {
				log.Printf("Retry attempt failed for %s: %v\n", j.path, err)
			}
		}(job)
	}
}

func translateFileRecursive(path string, excludeFolderName string) error {
	// read directory
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			if file.Name() == excludeFolderName {
				continue
			}

			err := translateFileRecursive(filepath.Join(path, file.Name()), excludeFolderName)
			if err != nil {
				log.Println(fmt.Errorf("failed to translate directory %s: %w", filepath.Join(path, file.Name()), err))
			}

			// Only translate markdown files
		} else if filepath.Ext(file.Name()) == ".md" {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				err := translateFile(filepath.Join(path, file.Name()), false)
				if err != nil {
					log.Println(fmt.Errorf("failed to translate file %s: %w", filepath.Join(path, file.Name()), err))
				}
			}()
			time.Sleep(3 * time.Second)
		}
	}

	return nil
}

func translateFile(path string, retrying bool) error {
	// read file all
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	defer func() {
		if err != nil && !retrying {
			retryQueue <- TranslationJob{
				path: path,
			}
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	prompt := `
제공될 문서를 모두 한국어로 깔끔하게 번역하고,다음과 같은 포맷을 지켜주세요.
- 모든 마크다운 형식은 html로 포팅해주세요.
- 순수 html 형식만 존재해야 합니다. (한국어 문서를 남겨두되 백틱 등은 존재하면 안됩니다.)
- 보기 쉽게 다크 테마를 적용해주세요. 포인트 컬러는 bluebird 컬러를 사용해주세요.
`

	// translate
	resp, err := client.GenerateText(prompt + string(content))
	if err != nil {
		return err
	}

	var response gemini.GeminiResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return err
	}

	if len(response.Candidates) == 0 {
		return fmt.Errorf("no candidates found")
	}

	if len(response.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("no candidates found")
	}

	// write file
	err = os.WriteFile(fmt.Sprintf("%s_translated.html", path), []byte(response.Candidates[0].Content.Parts[0].Text), 0645) // 0644 is the permission for the file
	if err != nil {
		return err
	}

	log.Printf("translated file %s\n", path)
	return nil
}

func main() {
	// read args
	args := os.Args
	if len(args) != 3 {
		fmt.Println("usage: docs-translater.exe <path> <exclude_folder_name>")
		return
	}

	// translate
	err := translateFileRecursive(args[1], args[2])
	if err != nil {
		log.Fatal(err)
	}

	// 모든 작업이 완료될 때까지 대기
	waitGroup.Wait()
	close(retryQueue) // 모든 작업이 완료되면 재시도 큐를 닫음
}
