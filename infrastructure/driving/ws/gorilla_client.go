package ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/josemontano1996/ai-chatbot-backend/pkg/utils"
)

type GorillaWSClient[T any] struct {
	Conn   *websocket.Conn
	Config WSConfig
}

func NewGorillaWSClient[T any](config WSConfig) (WSClientInterface[T], error) {
	err := utils.ValidateStruct(config)

	if err != nil {
		return nil, err
	}
	return &GorillaWSClient[T]{
		Config: config,
	}, nil
}

func (ws *GorillaWSClient[T]) Connect(ctx *gin.Context) error {
	conn, err := ws.upgradeConn(ctx)
	if err != nil {
		return fmt.Errorf("NewGorillaWSClient upgradeConn: %w", err)
	}

	conn.SetPingHandler(func(string) error {
		err := conn.WriteControl(
			websocket.PongMessage,
			[]byte{},
			time.Now().Add(10*time.Second),
		)
		return err
	})

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := conn.WriteControl(
				websocket.PingMessage,
				[]byte{},
				time.Now().Add(10*time.Second),
			); err != nil {
				return
			}
		}

	}()

	ws.Conn = conn
	return nil
}

func (client *GorillaWSClient[T]) ParseIncomingRequest() (payload *WSPayload[T], err error) {
	fmt.Println("parsing request", payload)
	err = client.Conn.ReadJSON(&payload)

	if err != nil {
		// Handle closure errors
		if websocket.IsUnexpectedCloseError(err,
			websocket.CloseNormalClosure,
			websocket.CloseGoingAway,
			websocket.CloseAbnormalClosure) {
			return nil, fmt.Errorf("websocket connection closed: %w", err)
		}
		return nil, fmt.Errorf("error reading JSON: %w", err)
	}

	err = utils.ValidateStruct(payload)
	if err != nil {
		return nil, fmt.Errorf("invalid message format: %w", err)
	}

	return
}

func (client *GorillaWSClient[T]) SendResposeToClient(response *WSPayload[T]) error {
	err := client.Conn.WriteJSON(response)
	if err != nil {
		return fmt.Errorf("error sending response to WS client: %w", err)
	}
	return nil
}

func (client *GorillaWSClient[T]) NewPayload(x T, err error) *WSPayload[T] {
	if err != nil {
		return &WSPayload[T]{
			Payload: x,
			Error:   err.Error(),
		}
	}
	return &WSPayload[T]{
		Payload: x,
	}

}

func (client *GorillaWSClient[T]) Disconnect() error {
	if client.Conn != nil {
		err := client.Conn.Close()
		if err != nil {
			return fmt.Errorf("error closing WS connection: %w", err)
		}
	}
	return nil
}

func (client *GorillaWSClient[T]) SendErrorToClient(err error) error {
	response := &WSPayload[T]{
		Payload: *new(T),
		Error:   err.Error(),
	}
	return client.SendResposeToClient(response)
}

func (ws *GorillaWSClient[T]) upgradeConn(ctx *gin.Context) (*websocket.Conn, error) {

	upgrader := websocket.Upgrader{
		ReadBufferSize: ws.Config.ReadBufferSize, WriteBufferSize: ws.Config.WriteBufferSize, CheckOrigin: func(r *http.Request) bool {
			return ws.Config.CheckOrigin()
		}}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade connection: %w", err)
	}

	return conn, nil
}
