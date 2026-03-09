package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Message is a simple chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMService wraps LLM chat and embedding calls
type LLMService interface {
	Chat(ctx context.Context, messages []Message) (string, error)
	Embed(ctx context.Context, text string) ([]float64, error)
}

type llmService struct {
	client *http.Client
	logger *zap.Logger
}

func NewLLMService(logger *zap.Logger) LLMService {
	return &llmService{
		client: &http.Client{Timeout: 60 * time.Second},
		logger: logger,
	}
}

func (s *llmService) Chat(ctx context.Context, messages []Message) (string, error) {
	apiKey := viper.GetString("ai.openai_api_key")
	baseURL := viper.GetString("ai.openai_base_url")
	model := viper.GetString("ai.chat_model")
	if apiKey == "" {
		return "", errors.New("AI API key 未配置，请设置 ai.openai_api_key")
	}

	reqBody := map[string]interface{}{
		"model":    model,
		"messages": messages,
	}
	data, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API 请求失败: %s", string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("LLM 返回空响应")
	}
	return result.Choices[0].Message.Content, nil
}

func (s *llmService) Embed(ctx context.Context, text string) ([]float64, error) {
	apiKey := viper.GetString("ai.openai_api_key")
	baseURL := viper.GetString("ai.openai_base_url")
	model := viper.GetString("ai.embedding_model")
	if apiKey == "" {
		return nil, errors.New("AI API key 未配置")
	}

	reqBody := map[string]interface{}{
		"model": model,
		"input": text,
	}
	data, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/embeddings", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Embedding API 请求失败: %s", string(body))
	}

	var result struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, errors.New("Embedding 返回空响应")
	}
	return result.Data[0].Embedding, nil
}
