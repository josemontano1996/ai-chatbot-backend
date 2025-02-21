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

var AllowedChatMessageCodes = []ChatMessageCode{UserChatMessageCode, AIBotChatMessageCode}

type ChatHistory []ChatMessage

type ChatMessage struct {
	UserId  string
	Code    ChatMessageCode
	Message string
}

func newMessage(code ChatMessageCode, message string, userId string) (*ChatMessage, error) {

	if !isValidMessageCode(code) {
		err := fmt.Errorf("invalid message code: %v, allowed codes: %v", code, AllowedChatMessageCodes)
		return nil, err
	}

	if code == UserChatMessageCode && userId == "" {
		err := fmt.Errorf("user id cannot be empty")
		return nil, err
	}

	if message == "" {
		err := fmt.Errorf("message cannot be empty")
		return nil, err
	}

	return &ChatMessage{
		Code:    code,
		Message: message,
		UserId:  userId,
	}, nil
}

func NewUserMessage(userId string, message string) (*ChatMessage, error) {
	return newMessage(UserChatMessageCode, message, userId)
}
func NewBotMessage(message string) (*ChatMessage, error) {
	return newMessage(AIBotChatMessageCode, message, "")
}
func NewSystemMessage(message string) (*ChatMessage, error) {
	return newMessage(AISystemChatMessageCode, message, "")
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
