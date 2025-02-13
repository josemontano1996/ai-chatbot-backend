package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//TODO: block unathorized domains, only allow our domain
		return true //Allowing all origins for now
	},
}

type WSClient struct {
	Conn *websocket.Conn
	//TODO: change type for uuid.UUID
	//UserId uuid.UUID
	UserId uuid.UUID
	// do not add the gin context as it would make it pretty heavy as it will persist
}

// Upgrades is responsible for taking the http request and
// upgrading it to a websocket connection
// The upgrader handles the initial websocket handshake
func upgradeConn(c *gin.Context) (*websocket.Conn, error) {
	return upgrader.Upgrade(c.Writer, c.Request, nil)
}

func NewWSClient(c *gin.Context, UserId uuid.UUID, timeout time.Duration) (*WSClient, error) {
	conn, err := upgradeConn(c)

	if err != nil {
		return nil, fmt.Errorf("error upgrading connection: %w", err)
	}

	// Configure Keep-Alive (Ping/Pong mechanism)
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

	// Set Write Wait to prevent dead connections
	conn.SetWriteDeadline(time.Now().Add(timeout))

	// Set ReadDeadline to enforce read timeouts
	conn.SetReadDeadline(time.Now().Add(timeout / 2))

	return &WSClient{
		Conn:   conn,
		UserId: UserId,
	}, nil
}
