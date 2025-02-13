package openai

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Docs for structured outputs: https://platform.openai.com/docs/guides/structured-outputs
type OpenAIConfig struct {
	UserId              uuid.UUID `json:"user_id" validate:"required"`
	APIKey              string    `json:"api_key" validate:"required"`
	Model               string    `json:"model" validate:"required"`
	MaxCompletionTokens uint32    `json:"max_completion_tokens" validate:"required,min=100,max=15000"`
	OptionalConfig      OptionalOpenAIConfig
}

type OptionalOpenAIConfig struct {
	Temperature     float32  `json:"temperature" validate:"omitempty"`
	PresencePenalty float32  `json:"presence_penalty" validate:"omitempty"`
	OutputType []string `json:"output_type" validate:"omitempty,dive,oneof=text json_object json_schema"`
}

var DefaultOptionalConfig = OptionalOpenAIConfig{
	OutputType:      []string{OpenAITextOutput},
	Temperature:     1.0,
	PresencePenalty: 0.0,
}

func newOpenAIConfig(userId uuid.UUID, apiKey string, model string, maxCompletionTokens uint32, optional OptionalOpenAIConfig) (*OpenAIConfig, error) {
	optionalConf := configureOptionalParams(optional)

	conf := &OpenAIConfig{
		UserId:              userId,
		APIKey:              apiKey,
		Model:               model,
		MaxCompletionTokens: maxCompletionTokens,
		OptionalConfig:      *optionalConf,
	}
	validate := validator.New()

	err := validate.Struct(conf)

	if err != nil {
		return nil, err
	}

	return conf, nil
}

func configureOptionalParams(params OptionalOpenAIConfig) *OptionalOpenAIConfig {
	optionalConf := DefaultOptionalConfig

	if len(params.OutputType) != 0 {
		optionalConf.OutputType = params.OutputType
	}

	if params.Temperature != 0 {
		optionalConf.Temperature = params.Temperature
	}
	if params.PresencePenalty != 0 {
		optionalConf.PresencePenalty = params.PresencePenalty
	}
	return &optionalConf
}
