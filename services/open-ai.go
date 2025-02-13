// package services

// import (
// 	"github.com/go-playground/validator/v10"
// 	"github.com/google/uuid"
// )

// // curl "https://api.openai.com/v1/chat/completions" \
// //     -H "Content-Type: application/json" \
// //     -H "Authorization: Bearer $OPENAI_API_KEY" \
// //     -d '{
// //         "model": "gpt-4o-mini",
// //         "messages": [
// //             {
// //                 "role": "system",
// //                 "content": "You are a helpful assistant."
// //             },
// //             {
// //                 "role": "user",
// //                 "content": "Write a haiku that explains the concept of recursion."
// //             }
// //         ]
// //     }'

// const (
// 	ModelOpenAIGpt4omini   string = "gpt-4o-mini"
// 	ModelOpenAIGpt4o       string = "gpt-4o"
// 	OpenAISystemRole       string = "system"
// 	OpenAIUserRole         string = "user"
// 	OpenAIBotRole          string = "assistant"
// 	OpenAITextOutput       string = "text"
// 	OpenAIJSONObjectOutput string = "json_object"
// 	OpenAIJSONSchemaOutput string = "json_schema"

// 	openAIUserRole   string = "user"
// 	openAISystemRole string = "developer"
// 	openAIBotRole    string = "assistant"
// )

// type OpenAIService struct {
// 	config *OpenAIConfig
// }

// type openAIMessage struct {
// 	Role    string `json:"role"`
// 	Content string
// }
// type openAIRequest struct {
// 	Model               string          `json:"model" validate:"required"`
// 	Messages            []openAIMessage `json:"messages" validate:"required"`
// 	MaxCompletionTokens uint32          `json:"max_completion_tokens" validate:"required"`
// 	Modalities          []string        `json:"modalities" validate:"required"`
// 	PresencePenalty     float32         `json:"presence_penalty" validate:"omitempty"`
// 	ResponseFormat      any             `json:"response_format" validate:"omitempty"`
// 	Temperature         float32         `json:"temperature" validate:"omitempty"`
// 	User                uuid.UUID       `json:"user" validate:"required"`
// }

// // Docs for structured outputs: https://platform.openai.com/docs/guides/structured-outputs
// type OpenAIConfig struct {
// 	UserId              uuid.UUID `json:"user_id" validate:"required"`
// 	APIKey              string    `json:"api_key" validate:"required"`
// 	Model               string    `json:"model" validate:"required"`
// 	MaxCompletionTokens uint32    `json:"max_completion_tokens" validate:"required,min=100,max=15000"`
// 	OptionalConfig      OptionalOpenAIConfig
// }

// type OptionalOpenAIConfig struct {
// 	Temperature     float32  `json:"temperature" validate:"omitempty"`
// 	PresencePenalty float32  `json:"presence_penalty" validate:"omitempty"`
// 	OutputType      []string `json:"output_type" validate:"omitempty,oneof=text json_object json_schema"`
// }

// var DefaultOptionalConfig = OptionalOpenAIConfig{
// 	OutputType:      []string{OpenAITextOutput},
// 	Temperature:     1.0,
// 	PresencePenalty: 0.0,
// }

// func NewOpenAIService(userId uuid.UUID, apiKey string, model string, maxCompletionTokens uint32, optional OptionalOpenAIConfig) (*OpenAIService, error) {
// 	conf, err := newOpenAIConfig(userId, apiKey, model, maxCompletionTokens, optional)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &OpenAIService{
// 		config: conf}, nil
// }

// func newOpenAIConfig(userId uuid.UUID, apiKey string, model string, maxCompletionTokens uint32, optional OptionalOpenAIConfig) (*OpenAIConfig, error) {
// 	optionalConf := newOptionalConfiguration(optional)

// 	conf := &OpenAIConfig{
// 		UserId:              userId,
// 		APIKey:              apiKey,
// 		Model:               model,
// 		MaxCompletionTokens: maxCompletionTokens,
// 		OptionalConfig:      *optionalConf,
// 	}
// 	validate := validator.New()

// 	err := validate.Struct(conf)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return conf, nil
// }

// func newOptionalConfiguration(params OptionalOpenAIConfig) *OptionalOpenAIConfig {
// 	optionalConf := DefaultOptionalConfig

// 	if len(params.OutputType) != 0 {
// 		optionalConf.OutputType = params.OutputType
// 	}

// 	if params.Temperature != 0 {
// 		optionalConf.Temperature = params.Temperature
// 	}
// 	if params.PresencePenalty != 0 {
// 		optionalConf.PresencePenalty = params.PresencePenalty
// 	}
// 	return &optionalConf
// }

// func CreateUserPrompt(prompt string) *openAIMessage {
// 	return &openAIMessage{
// 		Role:    openAIUserRole,
// 		Content: prompt,
// 	}
// }

// func CreateSystemPrompt(prompt string) *openAIMessage {
// 	return &openAIMessage{
// 		Role:    openAISystemRole,
// 		Content: prompt,
// 	}
// }

// func CreateBotPrompt(prompt string) *openAIMessage {
// 	return &openAIMessage{
// 		Role:    openAIBotRole,
// 		Content: prompt,
// 	}
// }
