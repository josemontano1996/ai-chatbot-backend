package entities

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ChatMessageCode int8

const (
	UserChatMessageCode  ChatMessageCode = 1
	AIBotChatMessageCode ChatMessageCode = 2
	// Chat messages should never display system information
	// AISystemChatMessageCode ChatMessageCode = 3
)

var AllowedChatMessageCodes = []ChatMessageCode{UserChatMessageCode, AIBotChatMessageCode}

type ChatMessage struct {
	Code    ChatMessageCode `json:"code" validate:"required,oneof=1 2"`
	Message string          `json:"message" binding:"required" validate:"required"`
}

func newMessage(code ChatMessageCode, message string) (*ChatMessage, error) {
	if !isValidMessageCode(code) {
		err := fmt.Errorf("invalid message code: %v, allowed codes: %v", code, AllowedChatMessageCodes)
		return nil, err
	}

	chatMessage := ChatMessage{
		Code:    code,
		Message: message}

	validator := validator.New()

	err := validator.Struct(chatMessage)

	if err != nil {
		err := fmt.Errorf("invalid message structure: %d", err)
		return nil, err
	}

	return &chatMessage, nil
}

func NewUserMessage(message string) (*ChatMessage, error) {
	return newMessage(UserChatMessageCode, message)
}
func NewBotMessage(message string) (*ChatMessage, error) {
	return newMessage(AIBotChatMessageCode, message)
}

func isValidMessageCode(code ChatMessageCode) bool {
	for _, allowedCode := range AllowedChatMessageCodes {
		if code == allowedCode {
			return true
		}
	}
	return false
}

// Chat history should never contain system messages
type ChatHistory []ChatMessage

func ParseArrayToChatHistory(array []string) (*ChatHistory, error) {
	messageHistory := make(ChatHistory, 0, len(array))

	for _, message := range array {
		var m ChatMessage
		err := json.Unmarshal([]byte(message), &m)

		if err != nil {
			return nil, fmt.Errorf("error parsing chat history array: %w", err)
		}

		messageHistory = append(messageHistory, m)

	}

	return &messageHistory, nil
}
