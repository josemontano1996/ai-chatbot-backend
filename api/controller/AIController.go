package controller

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/api/ws"
	"github.com/josemontano1996/ai-chatbot-backend/config"
	"github.com/josemontano1996/ai-chatbot-backend/repository"
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
	expirationTime := 5 * time.Minute

	client, err := ws.NewWSClient(c, userId, expirationTime)

	if err != nil {
		log.Fatal("fatal error connecting to websocket: ", err)
		return
	}

	defer client.Conn.Close()

	envs, err := config.LoadEnv("./", "app")

	if err != nil {
		log.Fatal("fatal error loading env variables: ", err)
		return
	}

	kv := repository.NewRedis(envs.RedisAddress, envs.RedisPassword, 0)

	// TODO change the fmt sprint to a string as it will be uuid string when db is up
	_, err = kv.Delete(c, "userkey").Result()

	if err != nil {
		log.Fatal("could not reset the kv value: ", err)
		return
	}

	log.Println("client connected: ", client)

	for {
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading json:", err)
			break
		}

		// 		El cliente envía el mensaje al servidor.

		// El servidor guarda el mensaje en Redis, junto con el resto del historial de la conversación.

		// El servidor recupera el historial de la conversación de Redis y lo usa como contexto para la IA.

		// La IA genera una respuesta.

		// El servidor envía la respuesta al cliente y la guarda en Redis.
		msgJSON, err := json.Marshal(msg)

		if err != nil {
			log.Println("Error marshaling message to JSON:", err)
			continue
		}

		err = kv.RPush(c, "userkey", msgJSON).Err()

		if err != nil {
			log.Println("Error pushing message to Redis:", err)
			continue
		}

		log.Printf("Received message: Type=%s, Content=%s\n", msg.Type, msg.Message)

		// Simulate Gemini API response - Replace with actual Gemini API call, now with context
		botResponse := simulateGeminiAPIResponse(msg.Message, conversations[client.UserId])

		responseMessage := Message{Type: "bot-message", Message: botResponse}

		responseMessageJson, err := json.Marshal(responseMessage)

		if err != nil {
			log.Println("Error Marshaling bot response to JSON", err)
			continue
		}

		err = kv.RPush(c, "userkey", responseMessageJson).Err()

		if err != nil {
			log.Println("Error pushing bot response to Redis:", err)
			continue
		}

		messages, err := kv.LRange(c, "userkey", 0, -1).Result() // Get the entire list
		if err != nil {
			log.Println("Error getting messages from Redis:", err)
			continue
		}

		var messageHistory []Message
		for _, message := range messages {
			var m Message
			err := json.Unmarshal([]byte(message), &m)
			if err != nil {
				log.Println("Error unmarshalling message:", err)
				continue // Skip the message if unmarshalling fails
			}
			messageHistory = append(messageHistory, m)
		}

		err = client.Conn.WriteJSON(messageHistory)
		if err != nil {
			log.Println("Error writing json:", err)
			break
		}
	}
}

func simulateGeminiAPIResponse(userMessage string, history History) string {
	//here i have to retrieve the message hisotory from redis
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
