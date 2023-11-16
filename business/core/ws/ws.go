package ws

import (
	"context"
	"fmt"

	"github.com/TSMC-Uber/server/business/sys/cachedb"
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
type Core struct{}

// NewCore constructs a core for user api access.
func NewCore() *Core {
	return &Core{}
}

func (c *Core) SendChatHistory(ctx context.Context, streamName string, channelName string, conn *websocket.Conn) error {
	messages, err := cachedb.XRange(ctx, streamName, "-", "+")
	if err != nil {
		return fmt.Errorf("xrange: %w", err)
	}

	for _, xMessage := range messages {
		msg, ok := xMessage.Values["message"].(string)
		if ok {
			conn.WriteMessage(websocket.TextMessage, []byte(msg))
		} else {
			return fmt.Errorf("message is not string: %w", err)
		}
	}

	return nil
}

func (c *Core) ReceiveChatMessages(ctx context.Context, streamName string, channelName string, conn *websocket.Conn, ch <-chan *redis.Message) error {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}
		// Add message to the Stream
		val, err := cachedb.XAdd(ctx, streamName, map[string]interface{}{"message": string(message)})
		if err != nil {
			return fmt.Errorf("xadd: %w", err)
		}
		fmt.Println("Message added to Stream:", val)
		// Publish message for real-time updates
		err = cachedb.Publish(ctx, channelName, string(message))
		if err != nil {
			return fmt.Errorf("publish: %w", err)
		}
	}
}
