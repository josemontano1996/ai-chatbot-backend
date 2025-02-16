package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//TODO: block unauthorized domains, only allow our domain
		return true // Allowing all origins for now
	},
}

type WSClient struct {
	Conn   *websocket.Conn
	UserID uuid.UUID // Consider using domain entity for User if relevant
}

type WSMessagePayload struct { // Define payload structure for WS messages
	Message string                `json:"message"`
	History *entities.ChatHistory `json:"history"`
}

func upgradeConn(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade connection: %w", err)
	}
	return conn, nil
}

func NewWSClient(c *gin.Context, userID uuid.UUID, timeout time.Duration) (*WSClient, error) {
	conn, err := upgradeConn(c)
	if err != nil {
		return nil, fmt.Errorf("NewWSClient upgradeConn: %w", err)
	}

	conn.SetPingHandler(func(string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(timeout))
		if err == websocket.ErrCloseSent {
			return err
		} else if err != nil {
			log.Println("Error sending pong:", err)
			return err
		}
		return nil
	})

	conn.SetPongHandler(func(string) error {
		return nil
	})

	conn.SetWriteDeadline(time.Now().Add(timeout))
	conn.SetReadDeadline(time.Now().Add(timeout / 2))

	return &WSClient{
		Conn:   conn,
		UserID: userID,
	}, nil
}

func (wsClient *WSClient) ParseIncomingRequest() (*WSMessagePayload, error) {
	var payload WSMessagePayload
	err := wsClient.Conn.ReadJSON(&payload)
	if err != nil {
		return nil, fmt.Errorf("ParseIncomingRequest ReadJSON: %w", err)
	}
	// You might want to add validation here for the payload structure

	return &payload, nil
}
