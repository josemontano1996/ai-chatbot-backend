package usecases

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/out"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
)

type AIChatUseCase[T any] struct {
	aiProvider  out.AIProvider[T]
	messageRepo outrepo.ChatMessageRepository
}

func NewAIChatUseCase[T any](aiProvider out.AIProvider[T], messageRepo outrepo.ChatMessageRepository) in.AIChatUseCase {
	return &AIChatUseCase[T]{
		aiProvider:  aiProvider,
		messageRepo: messageRepo,
	}
}

func (uc *AIChatUseCase[T]) SendChatMessage(ctx context.Context, userId string, userMessage string) (*in.AIChatResponse, error) {
	userMessageEntity, err := entities.NewUserMessage(userId, userMessage)
	if err != nil {
		return nil, err
	}

	//get history
	chatHistory, err := uc.messageRepo.GetChatHistory(ctx, userId)

	if err != nil {
		return nil, err
	}

	aiResponse, _, err := uc.aiProvider.SendMessage(ctx, userMessageEntity, chatHistory)

	if err != nil {
		return nil, err
	}

	// TODO: Implement logic to substract tokens from user's account
	// metadata.TokensSpent
	err = uc.messageRepo.SaveMessages(ctx, userId, userMessageEntity, aiResponse.ChatMessage)

	if err != nil {
		return nil, err
	}

	return aiResponse, nil
}
