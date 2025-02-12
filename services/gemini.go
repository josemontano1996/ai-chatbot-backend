package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"

	"google.golang.org/api/option"
)

const (
	// FinishReasonUnspecified means default value. This value is unused.
	FinishReasonUnspecified int8 = 0
	// FinishReasonStop means natural stop point of the model or provided stop sequence.
	FinishReasonStop int8 = 1
	// FinishReasonMaxTokens means the maximum number of tokens as specified in the request was reached.
	FinishReasonMaxTokens int8 = 2
	// FinishReasonSafety means the candidate content was flagged for safety reasons.
	FinishReasonSafety int8 = 3
	// FinishReasonRecitation means the candidate content was flagged for recitation reasons.
	FinishReasonRecitation int8 = 4
	// FinishReasonOther means unknown reason.
	FinishReasonOther int8 = 5
)

const (
	Gemini15FlashModelName string = "gemini-1.5-flash"
)

type GeminiService struct {
	ctx    context.Context
	client *genai.Client
	model  *genai.GenerativeModel
}

type ChatRequest struct {
	UserMessage sharedtypes.Message
	History     sharedtypes.History
}

type Prompt struct {
	UserPrompt string `json:"user_prompt" validate:"required"`
}

type AIResponse struct {
	Message    *sharedtypes.Message
	TokenCount int32
}

type AIServiceConfig struct {
	ModelName         string `json:"model_name" validate:"required"`
	SystemInstruction string `json:"system_instruction" validate:"required"`
	MaxOutputTokens   int32  `json:"max_output_tokens" validate:"required"`
	// ResponseMIMEType  string        `json:"response_mime_type,omitempty"`
	// ResponseSchema    *genai.Schema `json:"response_schema,omitempty"`
}

func NewAIServiceConfig(model string, systemInstruction string, maxOutputTokens int32) *AIServiceConfig {
	return &AIServiceConfig{
		ModelName:         model,
		SystemInstruction: systemInstruction,
		MaxOutputTokens:   maxOutputTokens,
	}
}

func NewGeminiService(ctx *gin.Context, apiKey string, config *AIServiceConfig) (*GeminiService, error) {

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	fmt.Println("connected to gemini: ", time.Now())
	if err != nil {
		return nil, err
	}
	model := client.GenerativeModel(config.ModelName)

	model.SystemInstruction = &genai.Content{
		Role:  "bot", // Or "model" depending on your requirements
		Parts: []genai.Part{genai.Text(config.SystemInstruction)},
	}
	model.MaxOutputTokens = &config.MaxOutputTokens
	// model.ResponseMIMEType = config.ResponseMIMEType
	// model.ResponseSchema = config.ResponseSchema
	var candidateNumber int32 = 1
	model.GenerationConfig.CandidateCount = &candidateNumber
	fmt.Println("finished setting up model config: ", time.Now())

	return &GeminiService{
		ctx:    ctx,
		client: client,
		model:  model}, nil
}

func (ai *GeminiService) Chat(ctx *gin.Context, userMessage *string, History *sharedtypes.History) (*AIResponse, error) {
	defer ai.client.Close()
	fmt.Println("started chat: ", time.Now())

	session := ai.model.StartChat()
	session.History = ai.parseHistory(History)
	fmt.Println("sent message ", time.Now())
	res, err := session.SendMessage(ctx, genai.Text(*userMessage))

	if err != nil {
		return nil, err
	}
	
	fmt.Println("parsing response: ", time.Now())
	parsedResponse, err := ai.parseAIRespose(res)
	if err != nil {
		return nil, err
	}

	return parsedResponse, nil
}
func (ai *GeminiService) parseHistory(History *sharedtypes.History) []*genai.Content {
	formattedHistory := make([]*genai.Content, 0)
	parsedHistoryElement := &genai.Content{}

	arrayLength := len(*History)
	for i, message := range *History {

		if i == arrayLength-1 && message.Type > 0 {
			// if the last message is from the message, then we break because we will give this message to the model throw a direct message.
			// if we do not break it will be added to the history and the model will be confused as the message is doubled
			break
		}

		parsedHistoryElement.Parts = append(parsedHistoryElement.Parts, genai.Text(message.Message))
		if message.Type < 0 {
			parsedHistoryElement.Role = "model"
		} else {
			parsedHistoryElement.Role = "user"
		}

		formattedHistory = append(formattedHistory, parsedHistoryElement)
	}

	return formattedHistory
}

func (ai *GeminiService) parseAIRespose(res *genai.GenerateContentResponse) (*AIResponse, error) {
	if len(res.Candidates) == 0 {
		return &AIResponse{}, errors.New("no candidates found")
	}

	candidate := res.Candidates[0]
	if candidate.FinishReason != 1 {
		_, err := ai.parseCandidateError(candidate)
		return &AIResponse{}, err
	}

	finalMessage := ""
	if cs := contentString(candidate.Content); cs != nil {
		finalMessage = *cs
	}

	return &AIResponse{
		Message: &sharedtypes.Message{
			Message: finalMessage,
			Type:    -1,
		},
		TokenCount: candidate.TokenCount,
	}, nil
}

// contentString converts genai.Content to a string. If the parts in the input content are of type
// text, they are concatenated with new lines in between them to form a string.
func contentString(c *genai.Content) *string {
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

func (ai *GeminiService) parseCandidateError(candidate *genai.Candidate) (int8, error) {
	switch candidate.FinishReason {
	case 2:
		return FinishReasonMaxTokens, errors.New("max tokens reached")
	case 3:
		return FinishReasonSafety, errors.New("safety reason")
	case 4:
		return FinishReasonRecitation, errors.New("recitation reason")
	case 5:
		return FinishReasonOther, errors.New("other reason")
	default:
		return FinishReasonUnspecified, nil
	}
}

// type GenerateContentConfig struct {
// 	// docs: https://cloud.google.com/vertex-ai/generative-ai/docs/model-reference/inference#generationconfig
// 	// https://cloud.google.com/vertex-ai/generative-ai/docs/multimodal/content-generation-parameters
// 	SystemInstruction string `json:"systemInstruction" validate:"required"`
// 	MaxOutputTokens   uint16 `json:"maxOutputTokens" validate:"required"`
// 	ResposeMimeType   string `json:"responseMimeType,omitempty"`
// 	ResponseSchema    string `json:"responseSchema,omitempty"`
// }

// func NewGeminiClient(c context.Context, apiKey string) (*GeminiClient, error) {
// 	client, err := genai.NewClient(c, option.WithAPIKey(apiKey))

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &GeminiClient{
// 		Client: client,
// 	}, nil
// }

// func (ai GeminiClient) GenerateContent(c *context.Context, model string, content []*genai.Content, config *genai.GenerateContentConfig) {
// 	ai.Client.Models.GenerateContent(*c, model, content, config)
// }

// func NewPrompt() {
// 	candidates := 1
// }
