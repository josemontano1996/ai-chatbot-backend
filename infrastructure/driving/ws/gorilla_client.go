package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type GorillaWSClient[T any] struct {
	Conn *websocket.Conn
}


func NewGorillaWSClient[T any]() WSClientInterface[T] {
	return &GorillaWSClient[T]{}
}

func (ws *GorillaWSClient[T]) Connect(config WSConfig) error {
	conn, err := upgradeConn(config)
	if err != nil {
		return fmt.Errorf("NewGorillaWSClient upgradeConn: %w", err)
	}
	
	conn.SetPingHandler(func(string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(config.ExpirationTime))
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
	
	conn.SetWriteDeadline(time.Now().Add(config.ExpirationTime))
	conn.SetReadDeadline(time.Now().Add(config.ExpirationTime / 2))
	
	ws.Conn = conn
	return nil
}

func (client *GorillaWSClient[T]) ParseIncomingRequest() (payload *WSPayload[T], err error) {
	err = client.Conn.ReadJSON(&payload)
	
	if err != nil {
		return nil, fmt.Errorf("error reading JSON from WS connection: %w", err)
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

func (client *GorillaWSClient[T]) NewPayload(x T) *WSPayload[T] {
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

func upgradeConn(config WSConfig) (*websocket.Conn, error) {

	upgrader := websocket.Upgrader{
		ReadBufferSize: config.ReadBufferSize, WriteBufferSize: config.WriteBufferSize, CheckOrigin: func(r *http.Request) bool {
			return config.CheckOrigin(r)
		}}

	conn, err := upgrader.Upgrade(config.Ctx.Writer, config.Ctx.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade connection: %w", err)
	}

	return conn, nil
}