package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/api/ws"
	"github.com/josemontano1996/ai-chatbot-backend/config"
	"github.com/josemontano1996/ai-chatbot-backend/handlers"
	"github.com/josemontano1996/ai-chatbot-backend/repository"
	"github.com/josemontano1996/ai-chatbot-backend/services"
	"github.com/josemontano1996/ai-chatbot-backend/services/openai"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

var userId uuid.UUID = uuid.New()

func PostAIController(c *gin.Context) {
	handleConnections(c) // Upgrade to websocket on POST request
}

func handleConnections(c *gin.Context) {
	expirationTime := 60 * time.Minute

	client, err := ws.NewWSClient(c, userId, expirationTime)

	if err != nil {
		log.Fatal("fatal error connecting to websocket: ", err)
		return
	}

	defer client.Conn.Close()
	fmt.Println("ws connectado: ", time.Now())

	envs, err := config.LoadEnv("./", "prod")
	if err != nil {
		log.Fatal("fatal error loading env variables: ", err)
		return
	}

	kv := repository.NewRedis(envs.RedisAddress, envs.RedisPassword, 0)


	// initialize the ai config
	optionalConfig := &openai.OptionalOpenAIConfig{}
	AIService, err := openai.NewOpenAIService(userId, envs.OpenAiApiKey, openai.ModelOpenAIGpt4omini, envs.MaxCompletionTokens, *optionalConfig)

	if err != nil {
		log.Fatal("could not create open ai config: ", err)
		return
	}
	//TODO: add a mechanism to block incoming requests if the model is processing
	// var isProcessing bool = false

	for {
		fmt.Println("inside the lopp: ", time.Now())

		userMessage, err := handlers.ParseUserMessageFromRequest(c)
		if err != nil {
			log.Println("Error parsing user message from request:", err)
			// TODO: handle parsing errors to the client via wesocket
		}

		msgHistory, err := kv.GetList(c, "userkey", 0, -1) // Get the entire list

		if err != nil {
			log.Println("Error getting messages from Redis:", err)
			continue
		}

		fmt.Println("boot open ai service: ", time.Now())
		gemini, err := services.NewGeminiService(c, envs.GeminiApiKey, aiConfig)

		if err != nil {
			// TODO: handle more gracefully
			log.Fatal("could not create gemini service: ", err)
			return
		}
		fmt.Println("parseando historial: ", time.Now())

		parsedMsgHistory := sharedtypes.ParseJSONToHistory(msgHistory)
		response, err := gemini.Chat(c, &userMessage.Message, parsedMsgHistory)

		if err != nil {
			//TODO: handle more gracefully
			log.Println("Error getting response from AI:", err)
			continue
		}

		//TODO: substract amount of tokens used

		responseMessageJSON, err := json.Marshal(response.Message)

		if err != nil {
			log.Println("Error Marshaling bot response to JSON", err)
			continue
		}

		fmt.Println("guarndando mensaje de bot en kv: ", time.Now())
		err = kv.RPush(c, "userkey", responseMessageJSON).Err()

		if err != nil {
			log.Println("Error pushing bot response to Redis:", err)
			continue
		}

		fmt.Println("cargando historial 2: ", time.Now())
		messages, err := kv.LRange(c, "userkey", 0, -1).Result() // Get the entire list
		if err != nil {
			log.Println("Error getting messages from Redis:", err)
			continue
		}

		updatedHistory := sharedtypes.NewHistory(messages)
		fmt.Println("historial parseado: ", time.Now())

		fmt.Println("enviando respesta a cliente: ", time.Now())

		err = client.Conn.WriteJSON(updatedHistory)
		if err != nil {
			log.Println("Error writing json:", err)
			fmt.Println("break: ", time.Now())

			break
		}
		fmt.Println("final: ", time.Now())

	}
}
