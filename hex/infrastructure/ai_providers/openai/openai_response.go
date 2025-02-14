package openaiadapter

import (
	"fmt"

	"github.com/josemontano1996/ai-chatbot-backend/hex/domain/entities"
	inputport "github.com/josemontano1996/ai-chatbot-backend/hex/domain/ports/input"
	outputport "github.com/josemontano1996/ai-chatbot-backend/hex/domain/ports/output"
)

// docs https://platform.openai.com/docs/api-reference/chat/object#chat/object-usage
type UsageObject struct {
	// Number of tokens in the prompt
	PromptTokens uint32 `json:"prompt_tokens" binding:"required"`
	// Number of tokens in the generated copmoletion
	CompletionTokens uint32 `json:"completion_tokens" binding:"required"`
	// Total number of tokens used in the request (prompt + completion)
	TotalTokens             uint32                  `json:"total_tokens" binding:"required"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details" binding:"nullable"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details" binding:"nullable"`
}

// PromptTokensDetails represents the nested "prompt_tokens_details" object.
type PromptTokensDetails struct {
	CachedTokens uint32 `json:"cached_tokens"`
}

// CompletionTokensDetails represents the nested "completion_tokens_details" object.
type CompletionTokensDetails struct {
	ReasoningTokens          uint32 `json:"reasoning_tokens"`
	AcceptedPredictionTokens uint32 `json:"accepted_prediction_tokens"`
	RejectedPredictionTokens uint32 `json:"rejected_prediction_tokens"`
}

type MessageObject struct {
	Index        uint8         `json:"index" binding:"required"`
	Message      MessageDetail `json:"message" binding:"required"`
	FinishReason string        `json:"finish_reason" binding:"required"`
}

type MessageDetail struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"nullable"`
	Refusal string `json:"refusal" binding:"nullable"`
}

type OpenAIResponse struct {
	ChatId    string          `json:"id" binding:"required"`
	CreatedAt int64           `json:"created" binding:"required"`
	Choices   []MessageObject `json:"choices" binding:"required"`
	Usage     UsageObject     `json:"usage" binding:"required"`
}

func (s *OpenAIAdapter) serializeResponse(data *OpenAIResponse) (*inputport.AIChatResponse, *outputport.AIResposeMetadata[*OpenAIResponse], error) { // Use inputport.ChatResponse
	totalTokensSpend := data.Usage.TotalTokens
	choice := data.Choices[0]
	message := choice.Message.Content

	if choice.Message.Refusal != "" {
		message = choice.Message.Refusal
	}

	msg, err := entities.NewMessage(entities.AIBotChatMessageCode, message)

	if err != nil {
		return &inputport.AIChatResponse{}, &outputport.AIResposeMetadata[*OpenAIResponse]{Metadata: data}, fmt.Errorf("error creating AI bot message: %w", err)
	}

	metadata := &outputport.AIResposeMetadata[*OpenAIResponse]{
		TokensSpent: totalTokensSpend,
		Metadata:    data,
	}
	
	return &inputport.AIChatResponse{
		ChatMessage: msg,
	}, metadata, nil
}