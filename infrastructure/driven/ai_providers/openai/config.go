package openaiadapter

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validOutputTypes = []string{OpenAITextOutput, OpenAIJSONObjectOutput, OpenAIJSONSchemaOutput}

// Docs for structured outputs: https://platform.openai.com/docs/guides/structured-outputs
type openAIConfig struct {
	UserId              uuid.UUID `json:"user_id" validate:"required"`
	APIKey              string    `json:"api_key" validate:"required"`
	Model               string    `json:"model" validate:"required"`
	MaxCompletionTokens uint32    `json:"max_completion_tokens" validate:"required,min=100,max=15000"`
	OptionalConfig      optionalOpenAIConfig
}

type optionalOpenAIConfig struct {
	Temperature     float32  `json:"temperature" validate:"omitempty"`
	PresencePenalty float32  `json:"presence_penalty" validate:"omitempty"`
	OutputType      []string `json:"output_type" validate:"omitempty,dive,oneof=text json_object json_schema"`
}

var DefaultOptionalConfig = optionalOpenAIConfig{
	OutputType:      []string{OpenAITextOutput},
	Temperature:     1.0,
	PresencePenalty: 0.0,
}

func NewOpenAIConfig(userId uuid.UUID, apiKey string, model string, maxCompletionTokens uint32, optionalConf optionalOpenAIConfig) (*openAIConfig, error) {

	conf := &openAIConfig{
		UserId:              userId,
		APIKey:              apiKey,
		Model:               model,
		MaxCompletionTokens: maxCompletionTokens,
		OptionalConfig:      optionalConf,
	}
	validate := validator.New()

	err := validate.Struct(conf)

	if err != nil {
		return nil, err
	}

	return conf, nil
}

func NewOpenAIOptionalConfig(temperature float32, presencePenalty float32, outputType []string) *optionalOpenAIConfig {
	optionalConf := DefaultOptionalConfig

	if isValidOutputType(outputType) {
		optionalConf.OutputType = outputType
	}

	if temperature > 0 {
		optionalConf.Temperature = temperature
	}
	if presencePenalty >= -2 && presencePenalty <= 2 {
		optionalConf.PresencePenalty = presencePenalty
	}

	return &optionalConf
}

func isValidOutputType(outputType []string) bool {
	if len(outputType) == 0 {
		return false
	}

	for _, t := range outputType {
		for _, validType := range validOutputTypes {
			if t == validType {
				return true
			}
		}
	}
	return false
}
