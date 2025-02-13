package gemini

type AIServiceConfig struct {
	ModelName         string `json:"model_name" validate:"required"`
	SystemInstruction string `json:"system_instruction" validate:"required"`
	MaxOutputTokens   int32  `json:"max_output_tokens" validate:"required"`
	//  ResponseMIMEType  string        `json:"response_mime_type,omitempty"`
	//  ResponseSchema    *genai.Schema `json:"response_schema,omitempty"`
}

func NewAIServiceConfig(model string, systemInstruction string, maxOutputTokens int32) *AIServiceConfig {
	return &AIServiceConfig{
		ModelName:         model,
		SystemInstruction: systemInstruction,
		MaxOutputTokens:   maxOutputTokens,
	}
}
