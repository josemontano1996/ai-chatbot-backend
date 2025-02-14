package openaiadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/josemontano1996/ai-chatbot-backend/domain/entities"
	inputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/input"
	outputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/output"
)

const (
	ModelOpenAIGpt4omini   string = "gpt-4o-mini"
	ModelOpenAIGpt4o       string = "gpt-4o"
	OpenAITextOutput       string = "text"
	OpenAIJSONObjectOutput string = "json_object"
	OpenAIJSONSchemaOutput string = "json_schema"
	openAIUserRole         string = "user"
	openAISystemRole       string = "developer"
	openAIBotRole          string = "assistant"
	openAIChatUrl          string = "https://api.openai.com/v1/chat/completions"
)

type OpenAIAdapter struct {
	config openAIConfig
	client *http.Client
}

type openAIMessage struct {
	Role    string `json:"role" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func NewOpenAIAdapter(config openAIConfig) *OpenAIAdapter {
	return &OpenAIAdapter{
		config: config,
		client: &http.Client{},
	}
}

func (ad *OpenAIAdapter) SendMessage(ctx context.Context, userMessage *entities.ChatMessage, prevHistory *entities.ChatHistory) (*inputport.AIChatResponse, *outputport.AIResposeMetadata[*OpenAIResponse], error) {
	instructions := "You are a helpful chatbot"
	prompts, err := ad.createPrompts(instructions, userMessage, prevHistory)

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[*OpenAIResponse]{}, err
	}

	req, err := ad.createRequestBody(prompts)

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[*OpenAIResponse]{}, err
	}

	resp, err := ad.client.Do(req)

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[*OpenAIResponse]{}, err
	}
	
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[*OpenAIResponse]{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[*OpenAIResponse]{}, errors.New("could not get response from OpenAI")
	}

	var openAIResponse OpenAIResponse
	err = json.Unmarshal(body, &openAIResponse)

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[*OpenAIResponse]{}, fmt.Errorf("could not unmarshal response %d ", err)
	}

	return ad.serializeResponse(&openAIResponse)
}

func (s *OpenAIAdapter) createRequestBody(prompts []openAIMessage) (*http.Request, error) {
	requestBody := openAIRequest{
		User:                s.config.UserId,
		Model:               s.config.Model,
		Messages:            prompts,
		MaxCompletionTokens: s.config.MaxCompletionTokens,
		Modalities:          s.config.OptionalConfig.OutputType,
		PresencePenalty:     s.config.OptionalConfig.PresencePenalty,
		Temperature:         s.config.OptionalConfig.Temperature,
	}

	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.New("could not marshal request body")
	}

	req, err := http.NewRequest("POST", openAIChatUrl, bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return nil, errors.New("could not create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	return req, nil
}

func (s *OpenAIAdapter) CloseConnection() {
	// No client to close in this simple HTTP client, but interface requires Close
}
