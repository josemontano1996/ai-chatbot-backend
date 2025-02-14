package usecases

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/domain/entities"
	inputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/input"
	outputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/output"
)

type AIChatUseCase[T any] struct {
	aiProvider  outputport.AIProvider[T]
	messageRepo outputport.ChatMessageRepository
}

func NewAIChatUseCase[T any](aiProvider outputport.AIProvider[T], messageRepo outputport.ChatMessageRepository) inputport.AIChatUseCase {
	return &AIChatUseCase[T]{
		aiProvider:  aiProvider,
		messageRepo: messageRepo,
	}
}

func (uc *AIChatUseCase[T]) SendChatMessage(ctx context.Context, user *entities.User, userMessage *entities.ChatMessage, history *entities.ChatHistory) (*inputport.AIChatResponse, error) {

	aiResponse, _, err := uc.aiProvider.SendMessage(ctx, userMessage, history)

	if err != nil {
		return &inputport.AIChatResponse{}, err
	}

	// TODO: Implement logic to substract tokens from user's account
	// metadata.TokensSpent
	err = uc.messageRepo.SaveMessages(ctx, user.ID.String(), userMessage, aiResponse.ChatMessage)

	if err != nil {
		return &inputport.AIChatResponse{}, err
	}

	return aiResponse, nil
}
