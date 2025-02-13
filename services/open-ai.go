package services

import "github.com/josemontano1996/ai-chatbot-backend/sharedtypes"

// curl "https://api.openai.com/v1/chat/completions" \
//     -H "Content-Type: application/json" \
//     -H "Authorization: Bearer $OPENAI_API_KEY" \
//     -d '{
//         "model": "gpt-4o-mini",
//         "messages": [
//             {
//                 "role": "system",
//                 "content": "You are a helpful assistant."
//             },
//             {
//                 "role": "user",
//                 "content": "Write a haiku that explains the concept of recursion."
//             }
//         ]
//     }'

const (
	ModelOpenAIGpt4omini   string = "gpt-4o-mini"
	ModelOpenAIGpt4o       string = "gpt-4o"
	OpenAISystemRole       string = "system"
	OpenAIUserRole         string = "user"
	OpenAIBotRole          string = "assistant"
	OpenAITextOutput       string = "text"
	OpenAIJSONObjectOutput string = "json_object"
	OpenAIJSONSchemaOutput string = "json_schema"
)

// Docs for structured outputs: https://platform.openai.com/docs/guides/structured-outputs
type OpenAIConfig struct {
	UserId              sharedtypes.UserId
	APIKey              string
	Model               string
	MaxTokens           uint32
	Temperature         float32
	PresencePenalty     float32
	MaxCompletionTokens uint32
	CandidatesNumber    uint8
	OutputType          string
}
