package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

const (
	ModelOpenAIGpt4omini   string = "gpt-4o-mini"
	ModelOpenAIGpt4o       string = "gpt-4o"
	OpenAISystemRole       string = "system"
	OpenAIUserRole         string = "user"
	OpenAIBotRole          string = "assistant"
	OpenAITextOutput       string = "text"
	OpenAIJSONObjectOutput string = "json_object"
	OpenAIJSONSchemaOutput string = "json_schema"

	openAIUserRole   string = "user"
	openAISystemRole string = "developer"
	openAIBotRole    string = "assistant"

	openAIChatUrl string = "https://api.openai.com/v1/chat/completions"
)

type OpenAIService struct {
	config *OpenAIConfig
}

type openAIMessage struct {
	Role    string `json:"role" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type openAIRequest struct {
	Model               string          `json:"model" validate:"required"`
	Messages            []openAIMessage `json:"messages" validate:"required"`
	MaxCompletionTokens uint32          `json:"max_completion_tokens" validate:"required"`
	Modalities          []string        `json:"modalities" validate:"required"`
	PresencePenalty     float32         `json:"presence_penalty" validate:"omitempty"`
	ResponseFormat      any             `json:"response_format" validate:"omitempty"`
	Temperature         float32         `json:"temperature" validate:"omitempty"`
	User                uuid.UUID       `json:"user" validate:"required"`
}




func NewOpenAIService(userId uuid.UUID, apiKey string, model string, maxCompletionTokens uint32, optional OptionalOpenAIConfig) (*OpenAIService, error) {
	conf, err := newOpenAIConfig(userId, apiKey, model, maxCompletionTokens, optional)

	if err != nil {
		return nil, err
	}

	return &OpenAIService{
		config: conf}, nil
}

func (s *OpenAIService) SendChatMessage(userMessage *sharedtypes.Message, prevHistory *[]openAIMessage) (*openAIRequest, error) {
	instructions := "You are a helpful chatbot"
	prompts := createPrompts(instructions, userMessage, prevHistory)
	req, err := s.createRequestBody(prompts)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("response status code was not 200: %d", resp.StatusCode)
			return nil, err
		}
		return nil, errors.New("could not send request")
	}

}

func (s *OpenAIService) createRequestBody(prompts *[]openAIMessage) (*http.Request, error) {

	requestBody := openAIRequest{
		User:                s.config.UserId,
		Model:               s.config.Model,
		Messages:            *prompts,
		MaxCompletionTokens: s.config.MaxCompletionTokens,
		Modalities:          s.config.OptionalConfig.OutputType,
		PresencePenalty:     s.config.OptionalConfig.PresencePenalty,
		Temperature:         s.config.OptionalConfig.Temperature,
	}

	jsonRequestBody, err := json.Marshal(requestBody)

	if err != nil {
		return nil, errors.New("could not marshal request body")
	}
	fmt.Println("creating request")
	req, err := http.NewRequest("POST", openAIChatUrl, bytes.NewBuffer(jsonRequestBody))

	if err != nil {
		return nil, errors.New("could not create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	return req, nil

}
