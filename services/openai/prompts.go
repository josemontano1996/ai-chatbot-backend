package openai

import "github.com/josemontano1996/ai-chatbot-backend/sharedtypes"

func createUserPrompt(prompt string) *openAIMessage {
	return &openAIMessage{
		Role:    openAIUserRole,
		Content: prompt,
	}
}

func createSystemPrompt(prompt string) *openAIMessage {
	return &openAIMessage{
		Role:    openAISystemRole,
		Content: prompt,
	}
}

func createBotPrompt(prompt string) *openAIMessage {
	return &openAIMessage{
		Role:    openAIBotRole,
		Content: prompt,
	}
}

func createPrompts(systemInstructions string, userMessage *sharedtypes.Message, prevHistory *sharedtypes.History) *[]openAIMessage {
	systemPrompt := createSystemPrompt(systemInstructions)
	userPrompt := createUserPrompt(userMessage.Message)
	mergedPrompts := mergePrompts(systemPrompt, userPrompt, prevHistory)
	return mergedPrompts
}

func mergePrompts(systemPrompt *openAIMessage, userPrompt *openAIMessage, history *sharedtypes.History) *[]openAIMessage {
	var mergedPrompts []openAIMessage

	// Create a new prompt with the system message
	mergedPrompts = append(mergedPrompts, *systemPrompt)
	// Add the history of prompts until now
	for _, prompt := range *history {
		mergedPrompts = append(mergedPrompts, *messageToOpenAIPrompt(prompt))
	}
	// Add the new user message prompt
	mergedPrompts = append(mergedPrompts, *userPrompt)

	return &mergedPrompts
}

func messageToOpenAIPrompt(message sharedtypes.Message) *openAIMessage {
	switch message.Code {
	case sharedtypes.AIBotResponseCode:
		return createBotPrompt(message.Message)
	case sharedtypes.AISystemMessageCode:
		return createSystemPrompt(message.Message)
	case sharedtypes.UserMessageCode:
		return createUserPrompt(message.Message)
	}
	return nil
}
