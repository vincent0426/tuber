package wsgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/core/ws"
	"github.com/TSMC-Uber/server/business/sys/cachedb"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handlers struct {
	ws   *ws.Core
	user *user.Core
}

// New constructs a handlers for route access.
func New(ws *ws.Core, user *user.Core) *Handlers {
	return &Handlers{
		ws:   ws,
		user: user,
	}
}

// Create adds a new trip to the system.
func (h *Handlers) Connect(ctx context.Context, c *gin.Context) error {
	userID := auth.GetUserID(c)
	// temp fix for testing, will remove later when auth is implemented
	// we need to implement a middleware to get user info from DB
	if userID == uuid.Nil {
		return fmt.Errorf("user not found")
	}

	// get user from DB
	user, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("query user by id: %w", err)
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return fmt.Errorf("upgrade: %w", err)
	}
	defer conn.Close()

	roomName := c.Request.URL.Query().Get("room")
	if roomName == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("No room specified"))
		return fmt.Errorf("no room specified")
	}

	streamName := "chatstream:" + roomName
	channelName := "chatroom:" + roomName

	// Send chat history to the new client from the Stream
	if err = h.ws.SendChatHistory(ctx, streamName, channelName, conn); err != nil {
		return fmt.Errorf("send chat history: %w", err)
	}

	// Subscribe to the Redis channel for real-time updates
	pubsub := cachedb.Subscribe(ctx, channelName)
	defer pubsub.Close()
	ch := pubsub.Channel()

	// Receive messages from WebSocket, send to Redis Stream, and publish to Redis channel
	go h.ws.ReceiveChatMessages(ctx, user, streamName, channelName, conn, ch)

	// Receive real-time messages from the Redis channel and send to WebSocket
	for msg := range ch {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			return fmt.Errorf("write message: %w", err)
		}
	}

	return nil
}
