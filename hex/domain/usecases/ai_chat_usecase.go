package usecases

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/hex/domain/entities"
	inputport "github.com/josemontano1996/ai-chatbot-backend/hex/domain/ports/input"
	outputport "github.com/josemontano1996/ai-chatbot-backend/hex/domain/ports/output"
)

type AIChatUseCase[T any] struct {
	aiProvider  outputport.AIProvider[T]
	messageRepo outputport.MessageRepository
	user        *entities.User
}

func NewChatUseCase[T any](user entities.User, aiProvider outputport.AIProvider[T], messageRepo outputport.MessageRepository) inputport.AIChatService {
	return &AIChatUseCase[T]{
		aiProvider:  aiProvider,
		messageRepo: messageRepo,
		user:        &user,
	}
}

func (uc *AIChatUseCase[T]) SendChatMessage(ctx context.Context, userMessage *entities.ChatMessage, history *entities.ChatHistory) (*inputport.AIChatResponse, error) {

	aiResponse, _, err := uc.aiProvider.SendMessage(ctx, userMessage, history)

	if err != nil {
		return &inputport.AIChatResponse{}, err
	}

	// TODO: Implement logic to substract tokens from user's account
	// metadata.TokensSpent
	err = uc.messageRepo.SaveMessages(ctx, uc.user.ID.String(), userMessage, aiResponse.ChatMessage)

	if err != nil {
		return &inputport.AIChatResponse{}, err
	}

	return aiResponse, nil
}
