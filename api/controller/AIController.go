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

	optionalConfig := openai.OptionalOpenAIConfig{}

	fmt.Println("creating open ai service")
	var AIService services.AIService[*openai.OpenAIResponse]
	AIService, err = openai.NewOpenAIService(userId, envs.OpenAiApiKey, openai.ModelOpenAIGpt4omini, envs.MaxCompletionTokens, &optionalConfig)

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

		jsonHistory, err := kv.GetList(c, "userkey", 0, -1) // Get the entire list

		if err != nil {
			log.Println("Error getting messages from Redis:", err)
			continue
		}

		prevHistory, err := sharedtypes.ParseJSONToHistory(jsonHistory)

		if err != nil {
			log.Println("error parsing json history to struct: ", err)
		}

		fmt.Println("send api call to open ai: ", time.Now())

		response, _, err := AIService.SendChatMessage(userMessage, prevHistory)

		fmt.Println("received open ai response: ", time.Now())

		if err != nil {
			log.Println("error when calling the open ai api: ", err)
			break
		}

		userMsgJson, err := json.Marshal(userMessage)

		if err != nil {
			log.Println("Error Marshaling bot response to JSON", err)
			continue
		}

		responseMessageJSON, err := json.Marshal(response.AIResponse)

		if err != nil {
			log.Println("Error Marshaling bot response to JSON", err)
			continue
		}

		// adding the user message and the ai response to the kv history
		_, err = kv.AddToList(c, "userkey", userMsgJson, responseMessageJSON)

		if err != nil {
			log.Println("Error pushing responses to Redis:", err)
			continue
		}

		err = client.Conn.WriteJSON(response.AIResponse)
		if err != nil {
			log.Println("Error writing json:", err)
			fmt.Println("break: ", time.Now())

			break
		}
		fmt.Println("final: ", time.Now())
		//TODO: substract amount of tokens used
		fmt.Println("reached bottom")
	}
}
