package controller

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/api/ws"
)

var userId int64 = 1

func PostAIController(c *gin.Context) {
	handleConnections(c) // Upgrade to websocket on POST request
}

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	// If an output structure is neede we can add it below
	// OutputStructure string `json:"output_structure"`
}

// History represents the list of messages in the session
type History []Message

// This will be stored and retrieved using redis
// TODO chang map tipe to uuid
var conversations = make(map[int64][]Message)

func handleConnections(c *gin.Context) {

	client, err := ws.NewWSClient(c, userId, 5*time.Minute)

	if err != nil {
		log.Fatal("fatal error connecting to websocket: ", err)
		return
	}

	defer client.Conn.Close()

	conversations[client.UserId] = []Message{}

	log.Println("client connected: ", client)

	for {
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading json:", err)
			break
		}

		log.Printf("Received message: Type=%s, Content=%s\n", msg.Type, msg.Message)

		conversations[client.UserId] = append(conversations[client.UserId], msg)
		// Append user message to history

		// Simulate Gemini API response - Replace with actual Gemini API call, now with context
		botResponse := simulateGeminiAPIResponse(msg.Message, conversations[client.UserId])

		responseMessage := Message{Type: "bot-message", Message: botResponse}
		conversations[client.UserId] = append(conversations[client.UserId], responseMessage) // Append bot response to history

		err = client.Conn.WriteJSON(conversations[client.UserId])
		if err != nil {
			log.Println("Error writing json:", err)
			break
		}
	}
}

func simulateGeminiAPIResponse(userMessage string, history History) string {

	context := ""
	if len(history) > 2 { // Example: Include last 2 messages as context
		context = "Previous messages: "
		for i := max(0, len(history)-3); i < len(history)-1; i++ { // Exclude the current user message
			if history[i].Type == "user-message" {
				context += "User: " + history[i].Message + "; "
			} else if history[i].Type == "bot-message" {
				context += "Bot: " + history[i].Message + "; "
			}
		}
	}

	return "AI Response: " + context + "You said: " + userMessage + ". This is a simulated response with context."
}
