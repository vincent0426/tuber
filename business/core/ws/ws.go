package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TSMC-Uber/server/business/sys/cachedb"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	// SendChatHistory(ctx context.Context, streamName string, channelName string, conn *websocket.Conn) error
	// Disconnect(ctx context.Context, userID uuid.UUID) error
}

// Core manages the set of APIs for user access.
type Core struct {
	storer Storer
}

// NewCore constructs a core for user api access.
func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) SendChatHistory(ctx context.Context, streamName string, channelName string, conn *websocket.Conn) error {
	messages, err := cachedb.XRange(ctx, streamName, "-", "+")
	if err != nil {
		return fmt.Errorf("xrange: %w", err)
	}

	for _, xMessage := range messages {
		jsonMsg, ok := xMessage.Values["message"].(string)
		if !ok {
			fmt.Println("Invalid message format in stream")
			continue
		}

		conn.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
	}

	return nil
}

func (c *Core) ReceiveChatMessages(ctx context.Context, userID uuid.UUID, streamName string, channelName string, conn *websocket.Conn, ch <-chan *redis.Message) error {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}

		chatMessage := ChatMessage{
			UserID:  userID.String(),
			Message: string(message),
		}

		jsonMessage, err := json.Marshal(chatMessage)
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}

		// Add message to the Stream
		val, err := cachedb.XAdd(ctx, streamName, map[string]interface{}{"message": jsonMessage})
		if err != nil {
			return fmt.Errorf("xadd: %w", err)
		}
		fmt.Println("Message added to Stream:", val)
		// Publish message for real-time updates
		err = cachedb.Publish(ctx, channelName, string(jsonMessage))
		if err != nil {
			return fmt.Errorf("publish: %w", err)
		}
	}
}
