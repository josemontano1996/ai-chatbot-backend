package sharedtypes

import (
	"encoding/json"
)

type UserId string

type Message struct {
	// 0 represents a message by the system
	// Greater than 0 represents a message by the user
	// Less than 0 represents a message by the bot
	Type    int8   `json:"type"`
	Message string `json:"message"`
	// If an output structure is neede we can add it below
	// OutputStructure string `json:"output_structure"`
}

func NewMessage(t int8, m string) Message {
	return Message{
		Type:    t,
		Message: m,
	}
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
