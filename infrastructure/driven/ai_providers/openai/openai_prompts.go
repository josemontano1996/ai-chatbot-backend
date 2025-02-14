package openaiadapter

import "github.com/josemontano1996/ai-chatbot-backend/domain/entities"


func (s *OpenAIAdapter) createPrompts(systemInstructions string, userMessage *entities.ChatMessage, prevHistory *entities.ChatHistory) ([]openAIMessage, error) {
	if userMessage == nil {
		userMessage = &entities.ChatMessage{}
	}

	systemPrompt := s.createSystemPrompt(systemInstructions)
	userPrompt := s.createUserPrompt(userMessage.Message)

	var history entities.ChatHistory
	if prevHistory != nil {
		history = *prevHistory
	}

	merged := s.mergePrompts(systemPrompt, userPrompt, history)
	return merged, nil
}

func (s *OpenAIAdapter) mergePrompts(systemPrompt openAIMessage, userPrompt openAIMessage, history entities.ChatHistory) []openAIMessage {
	var mergedPrompts []openAIMessage
	mergedPrompts = append(mergedPrompts, systemPrompt)

	for _, msg := range history {
		prompt := s.messageToOpenAIPrompt(msg)
		if prompt != nil {
			mergedPrompts = append(mergedPrompts, *prompt)
		}
	}

	mergedPrompts = append(mergedPrompts, userPrompt)
	return mergedPrompts
}

func (s *OpenAIAdapter) messageToOpenAIPrompt(message entities.ChatMessage) *openAIMessage {
	if message.Message == "" {
		return nil
	}
	switch message.Code {
	case entities.AIBotChatMessageCode:
		p := s.createBotPrompt(message.Message)
		return &p
	case entities.UserChatMessageCode:
		p := s.createUserPrompt(message.Message)
		return &p
	default:
		p := s.createSystemPrompt(message.Message)
		return &p
	}
}

func (s *OpenAIAdapter) createUserPrompt(prompt string) openAIMessage {
	return openAIMessage{
		Role:    openAIUserRole,
		Content: prompt,
	}
}

func (s *OpenAIAdapter) createSystemPrompt(prompt string) openAIMessage {
	return openAIMessage{
		Role:    openAISystemRole,
		Content: prompt,
	}
}

func (s *OpenAIAdapter) createBotPrompt(prompt string) openAIMessage {
	return openAIMessage{
		Role:    openAIBotRole,
		Content: prompt,
	}
}
