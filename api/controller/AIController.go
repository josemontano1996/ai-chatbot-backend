package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func PostAIController(c *gin.Context) {
	handleConnections(c) // Upgrade to websocket on POST request
}

// Upgrades is responsible for taking the http request and
// upgrading it to a websocket connection
// The upgrader handles the initial websocket handshake
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//TODO: block unathorized domains, only allow our domain
		return true //Allowing all origins for now
	},
}

type Client struct {
	conn *websocket.Conn
	//TODO: change type for uuid.UUID
	//UserId uuid.UUID
	UserId int64
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// This will be stored and retrieved using redis
// TODO chang map tipe to uuid
var conversations = make(map[int64][]Message)

func handleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Fatal("error upgrading connection: ", err)
		return
	}

	defer conn.Close()

	client := &Client{conn: conn,
		UserId: 1}

	conversations[client.UserId] = []Message{}

	log.Println("client connected: ", client)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading json:", err)
			break
		}

		log.Printf("Received message: Type=%s, Content=%s\n", msg.Type, msg.Content)

		conversations[client.UserId] = append(conversations[client.UserId], msg)
		// Append user message to history

		// Simulate Gemini API response - Replace with actual Gemini API call, now with context
		botResponse := simulateGeminiAPIResponse(msg.Content, conversations[client.UserId])

		responseMessage := Message{Type: "bot-message", Content: botResponse}
		conversations[client.UserId] = append(conversations[client.UserId], responseMessage) // Append bot response to history

		err = conn.WriteJSON(responseMessage)
		if err != nil {
			log.Println("Error writing json:", err)
			break
		}
	}
}

func simulateGeminiAPIResponse(userMessage string, history []Message) string {

	context := ""
	if len(history) > 2 { // Example: Include last 2 messages as context
		context = "Previous messages: "
		for i := max(0, len(history)-3); i < len(history)-1; i++ { // Exclude the current user message
			if history[i].Type == "user-message" {
				context += "User: " + history[i].Content + "; "
			} else if history[i].Type == "bot-message" {
				context += "Bot: " + history[i].Content + "; "
			}
		}
	}

	return "AI Response: " + context + "You said: " + userMessage + ". This is a simulated response with context."
}
