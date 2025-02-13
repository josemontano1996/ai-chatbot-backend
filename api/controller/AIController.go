package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/api/ws"
	"github.com/josemontano1996/ai-chatbot-backend/config"
	"github.com/josemontano1996/ai-chatbot-backend/repository"
	"github.com/josemontano1996/ai-chatbot-backend/services"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

var userId int64 = 1

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

	// TODO change the fmt sprint to a string as it will be uuid string when db is up
	redisHistory, err := kv.LRange(c, "userkey", 0, -1).Result()

	if err != nil {
		log.Fatal("could not reset the kv value: ", err)
		return
	}

	prevHistory, err := sharedtypes.ParseJSONToHistory(redisHistory)

	if err != nil {
		fmt.Println("error parsing history: ", err)
	}

	if len(*prevHistory) == 0 {
		fmt.Println("prev history: ", prevHistory)
		//TODO: inject system prompt at the beginnig of the history
	} else {
		//Send previous history to the client
		err = client.Conn.WriteJSON(prevHistory)
		if err != nil {
			log.Println("Error writing json: ", err)
		}
	}

	// initialize the ai config
	

	for {
		fmt.Println("inside the lopp: ", time.Now())

		var userMessage sharedtypes.Message
		fmt.Println("leyendo json de la petiicon ", time.Now())

		err := client.Conn.ReadJSON(&userMessage)
		if err != nil {
			log.Println("Error reading json:", err)
			break
		}
		userMessage.Type = 1

		msgHistory, err := kv.LRange(c, "userkey", 0, -1).Result() // Get the entire list

		if err != nil {
			log.Println("Error getting messages from Redis:", err)
			continue
		}
		fmt.Println("boot gemini service: ", time.Now())

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
