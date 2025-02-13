package sharedtypes

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type UserId uuid.UUID

const (
	UserMessageCode     int8 = 1
	AIBotResponseCode   int8 = 2
	AISystemMessageCode int8 = 3
)

var AllowedMessageCodes = []int8{UserMessageCode, AIBotResponseCode, AISystemMessageCode}

type Message struct {
	// 0 represents a message by the system
	// Greater than 0 represents a message by the user
	// Less than 0 represents a message by the bot
	Code    int8   `json:"code"`
	Message string `json:"message" validate:"required"`
	// If an output structure is neede we can add it below
	// OutputStructure string `json:"output_structure"`
}

func NewMessage(code int8, msg string) (*Message, error) {
	// Validate the message code
	validCode := false

	for _, allowedCode := range AllowedMessageCodes {
		if code == allowedCode {
			validCode = true
			break
		}
	}

	if !validCode {
		err := fmt.Errorf("invalid message code: %d", code)
		return &Message{}, err
	}

	return &Message{
		Code:    code,
		Message: msg,
	}, nil
}

// History represents the list of messages in the session
type History []Message

func ParseJSONToHistory(messages []string) (*History, error) {
	var messageHistory History
	for _, message := range messages {
		var m Message
		err := json.Unmarshal([]byte(message), &m)
		if err != nil {
			return nil, err
		}
		messageHistory = append(messageHistory, m)
	}
	return &messageHistory, nil
}
