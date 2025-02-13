package openai

type UsageObject struct {
	PromptTokens            int                     `json:"prompt_tokens" binding:"required"`
	CompletionTokens        int                     `json:"completion_tokens" binding:"required"`
	TotalTokens             int                     `json:"total_tokens" binding:"required"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details" binding:"nullable"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details" binding:"nullable"`
}

// PromptTokensDetails represents the nested "prompt_tokens_details" object.
type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

// CompletionTokensDetails represents the nested "completion_tokens_details" object.
type CompletionTokensDetails struct {
	ReasoningTokens          int `json:"reasoning_tokens"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
}

type MessageObject struct {
	Index        int           `json:"index" binding:"required"`
	Message      MessageDetail `json:"message" binding:"required"`
	FinishReason string        `json:"finish_reason" binding:"required"`
}

type MessageDetail struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"nullable"`
	Refusal string `json:"refusal" binding:"nullable"`
}

type OpenAIResponse struct {
	ChatId    string        `json:"id" binding:"required"`
	CreatedAt int64         `json:"created" binding:"required"`
	Choices   MessageObject `json:"choices" binding:"required"`
	Usage     UsageObject   `json:"usage" binding:"required"`
}
