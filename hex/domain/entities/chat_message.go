package entities

import (
	"encoding/json"
	"fmt"
)

type ChatMessageCode int8

const (
	UserChatMessageCode     ChatMessageCode = 1
	AIBotChatMessageCode    ChatMessageCode = 2
	AISystemChatMessageCode ChatMessageCode = 3
)

var AllowedChatMessageCodes = []ChatMessageCode{UserChatMessageCode, AIBotChatMessageCode, AISystemChatMessageCode}

type ChatMessage struct {
	Code    ChatMessageCode `json:"code"`
	Message string          `json:"message" binding:"required" validate:"required"`
}

func NewMessage(code ChatMessageCode, message string) (*ChatMessage, error) {
	// Validate the message code
	validCode := false

	for _, allowedCode := range AllowedChatMessageCodes {
		if code == allowedCode {
			validCode = true
			break
		}
	}

	if !validCode {
		err := fmt.Errorf("invalid message code: %v, allowed codes: %v", code, AllowedChatMessageCodes)
		return &ChatMessage{}, err
	}

	return &ChatMessage{
		Code:    code,
		Message: message,
	}, nil
}

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
