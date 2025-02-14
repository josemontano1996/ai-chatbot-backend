package openaiadapter

import "github.com/google/uuid"

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