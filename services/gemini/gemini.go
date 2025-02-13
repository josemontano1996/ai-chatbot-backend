package gemini

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/josemontano1996/ai-chatbot-backend/services"
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
	GeminiTextOutput       string = "text"

	geminiUserRole   string = "user"
	geminiSystemRole string = "system"
	geminiBotRole    string = "model"
)

type GeminiService struct {
	ctx    *gin.Context
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

func NewGeminiService(ctx *gin.Context, apiKey string, config *AIServiceConfig) (*GeminiService, error) {

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		return nil, err
	}
	model := client.GenerativeModel(config.ModelName)

	model.SystemInstruction = &genai.Content{
		Role:  geminiSystemRole,
		Parts: []genai.Part{genai.Text(config.SystemInstruction)},
	}
	model.MaxOutputTokens = &config.MaxOutputTokens
	// model.ResponseMIMEType = config.ResponseMIMEType
	// model.ResponseSchema = config.ResponseSchema
	var candidateNumber int32 = 1
	model.GenerationConfig.CandidateCount = &candidateNumber

	return &GeminiService{
		ctx:    ctx,
		client: client,
		model:  model}, nil
}

func (s GeminiService) SendChatMessage(userMessage *sharedtypes.Message, prevHistory *sharedtypes.History) (response *services.ChatResponse, metadata any, err error) {
	defer s.client.Close()

	chatSession := s.model.StartChat()

	chatSession.History = s.parseHistory(prevHistory)

	geminiResponse, err := chatSession.SendMessage(s.ctx, genai.Text(userMessage.Message))

	if err != nil {
		return
	}

	response, err = s.serializeResponse(geminiResponse)
	metadata = ""
	return
}

func (ai *GeminiService) parseHistory(History *sharedtypes.History) []*genai.Content {
	formattedHistory := make([]*genai.Content, 0)
	parsedHistoryElement := &genai.Content{}

	arrayLength := len(*History)
	for i, message := range *History {

		if i == arrayLength-1 && message.Code == sharedtypes.UserMessageCode {
			// if the last message is from the message, then we break because we will give this message to the model throw a direct message.
			// if we do not break it will be added to the history and the model will be confused as the message is doubled
			break
		}

		parsedHistoryElement.Parts = append(parsedHistoryElement.Parts, genai.Text(message.Message))
		if message.Code == sharedtypes.UserMessageCode {
			parsedHistoryElement.Role = geminiUserRole
		} else {
			parsedHistoryElement.Role = geminiBotRole
		}

		formattedHistory = append(formattedHistory, parsedHistoryElement)
	}

	return formattedHistory
}

func (ai *GeminiService) serializeResponse(res *genai.GenerateContentResponse) (*services.ChatResponse, error) {
	if len(res.Candidates) == 0 {
		return &services.ChatResponse{}, errors.New("no candidates found")
	}

	candidate := res.Candidates[0]
	if candidate.FinishReason != 1 {
		_, err := ai.parseCandidateError(candidate)
		return &services.ChatResponse{}, err
	}

	finalMessage := ""
	if cs := contentString(candidate.Content); cs != nil {
		finalMessage = *cs
	}

	return &services.ChatResponse{
		AIResponse: &sharedtypes.Message{
			Code:    sharedtypes.AIBotResponseCode,
			Message: finalMessage,
		},
		TotalTokensSpend: uint32(candidate.TokenCount),
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
