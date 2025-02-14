package geminiadapter

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"github.com/josemontano1996/ai-chatbot-backend/domain/entities"
	inputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/input"
	outputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/output"
	"google.golang.org/api/option"
)

const (
	Gemini15FlashModelName string = "gemini-1.5-flash"
	GeminiTextOutput       string = "text"
	geminiUserRole         string = "user"
	geminiSystemRole       string = "system"
	geminiBotRole          string = "model"
)

var defaulModelName = Gemini15FlashModelName
var validModelNames = []string{Gemini15FlashModelName}

type GeminiAdapter struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

var minimumOutputTokens int32 = 500

type geminiConfig struct { //  Config struct for adapter initialization
	APIKey          string `validate:"required"`
	ModelName       string `validate:"required"`
	MaxOutputTokens int32  `validate:"required,min=500"`
	candidateNumber int32  `validate:"required,min=1"`
}

func NewGeminiAdapter(ctx context.Context, config geminiConfig) (*GeminiAdapter, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))

	if err != nil {
		return &GeminiAdapter{}, err
	}

	model := client.GenerativeModel(config.ModelName)

	model.SetMaxOutputTokens(config.MaxOutputTokens)
	model.SetCandidateCount(config.candidateNumber)

	return &GeminiAdapter{
		client: client,
		model:  model,
	}, nil
}

// Creates a validated GeminiConfig struct
func NewGeminiConfig(apiKey string, modelName string, maxOutputTokens int32) (*geminiConfig, error) {
	configStruct := geminiConfig{}

	for _, model := range validModelNames {
		if model == modelName {
			configStruct.ModelName = modelName
			break
		}
	}

	if configStruct.ModelName == "" {
		configStruct.ModelName = defaulModelName
	}

	configStruct.APIKey = apiKey

	configStruct.candidateNumber = 1

	if maxOutputTokens < 500 {
		configStruct.MaxOutputTokens = minimumOutputTokens
	} else {
		configStruct.MaxOutputTokens = maxOutputTokens
	}

	validator := validator.New()

	err := validator.Struct(configStruct)

	if err != nil {
		return &geminiConfig{}, err
	}

	return &configStruct, nil
}

func (ad *GeminiAdapter) SendMessage(ctx context.Context, userMessage *entities.ChatMessage, prevHistory *entities.ChatHistory) (*inputport.AIChatResponse, *outputport.AIResposeMetadata[any], error) {

	ad.model.SystemInstruction = &genai.Content{
		Role:  geminiSystemRole,
		Parts: []genai.Part{genai.Text("You are a haleful assistant, help the user with their enquiries")},
	}

	chatSession := ad.model.StartChat()
	chatSession.History = ad.parseChatHistory(prevHistory)

	geminiResponse, err := chatSession.SendMessage(ctx, genai.Text(userMessage.Message))

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[any]{}, err
	}

	aiChatResponse, metadata, err := ad.serializeResponse(geminiResponse)

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[any]{}, err
	}

	return aiChatResponse, metadata, nil
}

func (ad *GeminiAdapter) parseChatHistory(chatHistory *entities.ChatHistory) []*genai.Content {
	formattedChatHistory := make([]*genai.Content, 0)

	// Avoiding nil pointer exception
	if chatHistory == nil {
		return formattedChatHistory
	}

	parsedHistoryElement := &genai.Content{}
	arrayLength := len(*chatHistory)

	for i, message := range *chatHistory {
		if i == arrayLength-1 && message.Code == entities.UserChatMessageCode {
			break // Donot add last user message to the history
		}
		parsedHistoryElement.Parts = append(parsedHistoryElement.Parts, genai.Text(message.Message))

		if message.Code == entities.UserChatMessageCode {
			parsedHistoryElement.Role = geminiUserRole
		} else {
			parsedHistoryElement.Role = geminiBotRole
		}
	}

	return formattedChatHistory
}

func (ad *GeminiAdapter) serializeResponse(res *genai.GenerateContentResponse) (*inputport.AIChatResponse, *outputport.AIResposeMetadata[any], error) {

	if len(res.Candidates) == 0 {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[any]{}, errors.New("no candidates found in response")
	}

	candidate := res.Candidates[0]

	if candidate.FinishReason != genai.FinishReasonStop {
		_, err := ad.parseCandidateError(candidate)
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[any]{}, err
	}

	finalMessage := ""

	if cs := contentToString(candidate.Content); cs != nil {
		finalMessage = *cs
	}

	msg, err := entities.NewBotMessage(finalMessage)

	if err != nil {
		return nil, &outputport.AIResposeMetadata[any]{}, fmt.Errorf("error creating AI bot message: %w", err)
	}

	return &inputport.AIChatResponse{
			ChatMessage: msg,
		}, &outputport.AIResposeMetadata[any]{
			TokensSpent: uint32(candidate.TokenCount),
		}, nil
}

func (ad *GeminiAdapter) parseCandidateError(candidate *genai.Candidate) (int8, error) {
	switch candidate.FinishReason {
	case genai.FinishReasonMaxTokens:
		return 2, errors.New("max tokens reached")
	case genai.FinishReasonSafety:
		return 3, errors.New("safety reason")
	case genai.FinishReasonRecitation:
		return 4, errors.New("recitation reason")
	case genai.FinishReasonOther:
		return 5, errors.New("other reason")
	default:
		return 0, fmt.Errorf("unspecified finish reason: %v", candidate.FinishReason)
	}
}

func contentToString(c *genai.Content) *string {
	if c == nil || c.Parts == nil {
		return nil
	}

	cStrs := make([]string, len(c.Parts))
	for i, part := range c.Parts {
		if pt, ok := part.(genai.Text); ok {
			cStrs[i] = string(pt)
		} else {
			return nil
		}
	}

	cStr := strings.Join(cStrs, "\n")
	return &cStr
}

func (ad *GeminiAdapter) CloseConnection() {
	if ad.client != nil {
		ad.client.Close()
	}
}
