package openai

import (
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

// Changed return type to value type for safety
func createUserPrompt(prompt string) openAIMessage {
	return openAIMessage{
		Role:    openAIUserRole,
		Content: prompt,
	}
}

// Changed return type to value type
func createSystemPrompt(prompt string) openAIMessage {
	return openAIMessage{
		Role:    openAISystemRole,
		Content: prompt,
	}
}

// Changed return type to value type
func createBotPrompt(prompt string) openAIMessage {
	return openAIMessage{
		Role:    openAIBotRole,
		Content: prompt,
	}
}

// Changed return type to value slice and added error handling
func createPrompts(systemInstructions string, userMessage *sharedtypes.Message, prevHistory *sharedtypes.History) ([]openAIMessage, error) {
	if userMessage == nil {
		userMessage = &sharedtypes.Message{}
	}

	systemPrompt := createSystemPrompt(systemInstructions)
	userPrompt := createUserPrompt(userMessage.Message)

	// Handle nil history
	var history []sharedtypes.Message
	if prevHistory != nil {
		history = *prevHistory
	}

	merged := mergePrompts(systemPrompt, userPrompt, history)
	return merged, nil
}

// Changed parameters to value types and added safety checks
func mergePrompts(systemPrompt openAIMessage, userPrompt openAIMessage, history []sharedtypes.Message) []openAIMessage {
	var mergedPrompts []openAIMessage

	// Add system prompt (always present)
	mergedPrompts = append(mergedPrompts, systemPrompt)

	// Process history safely
	for _, msg := range history {
		prompt := messageToOpenAIPrompt(msg)
		if prompt != nil {
			mergedPrompts = append(mergedPrompts, *prompt)
		}
	}

	// Add user prompt (always present)
	mergedPrompts = append(mergedPrompts, userPrompt)

	return mergedPrompts
}

// Changed parameter to value type and added validation
func messageToOpenAIPrompt(message sharedtypes.Message) *openAIMessage {
	if message.Message == "" {
		return nil
	}

	switch message.Code {
	case sharedtypes.AIBotResponseCode:
		p := createBotPrompt(message.Message)
		return &p
	case sharedtypes.AISystemMessageCode:
		p := createSystemPrompt(message.Message)
		return &p
	case sharedtypes.UserMessageCode:
		p := createUserPrompt(message.Message)
		return &p
	default:
		return nil
	}
}
