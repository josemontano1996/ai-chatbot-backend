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

func NewMessage(code ChatMessageCode, message string) (*ChatMessage, error) {
	if !isValidMessageCode(code) {
		err := fmt.Errorf("invalid message code: %v, allowed codes: %v", code, AllowedChatMessageCodes)
		return &ChatMessage{}, err
	}

	chatMessage := ChatMessage{
		Code:    code,
		Message: message}

	validator := validator.New()

	err := validator.Struct(chatMessage)

	if err != nil {
		err := fmt.Errorf("invalid message structure: %d", err)
		return &ChatMessage{}, err
	}

	return &chatMessage, nil
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

func ParseArrayToChatHistory(messageArray []string) (*ChatHistory, error) {
	var messageHistory ChatHistory

	var m ChatMessage
	for _, message := range messageArray {
		err := json.Unmarshal([]byte(message), &m)

		if err != nil {
			err := fmt.Errorf("error parsing chat history array: %d", err)
			return &ChatHistory{}, err
		}

		messageHistory = append(messageHistory, m)

	}

	return &messageHistory, nil
}
