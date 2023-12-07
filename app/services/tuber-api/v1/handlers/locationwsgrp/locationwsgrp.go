package locationwsgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/locationws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handlers struct {
	locationws *locationws.Core
}

// New constructs a handlers for route access.
func New(locationws *locationws.Core) *Handlers {
	return &Handlers{
		locationws: locationws,
	}
}

var upgrader = websocket.Upgrader{
	// TODO: check origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handlers) DriverWebSocketHandler(ctx context.Context, c *gin.Context) error {
	// userID := auth.GetUserID(c)
	// temp fix for testing, will remove later when auth is implemented
	// we need to implement a middleware to get user info from DB
	// if userID == uuid.Nil {
	// 	id, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	// 	if err != nil {
	// 		return fmt.Errorf("parse uuid: %w", err)
	// 	}

	// 	userID = id
	// }

	// upgrade get request to websocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return fmt.Errorf("upgrade: %w", err)
	}
	defer conn.Close()

	tripID := c.Request.URL.Query().Get("trip_id")
	if tripID == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("No trip specified"))
		return fmt.Errorf("no trip specified")
	}

	// TODO: check if user is driver for this trip

	_, err = locationws.ServeClient(
		conn,
		tripID,
		true,
	)
	if err != nil {
		return fmt.Errorf("serve client: %w", err)
	}
	return nil
}

func (h *Handlers) PassengerWebSocketHandler(ctx context.Context, c *gin.Context) error {
	// userID := auth.GetUserID(c)
	// temp fix for testing, will remove later when auth is implemented
	// we need to implement a middleware to get user info from DB
	// if userID == uuid.Nil {
	// 	id, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	// 	if err != nil {
	// 		return fmt.Errorf("parse uuid: %w", err)
	// 	}

	// 	userID = id
	// }

	// upgrade get request to websocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return fmt.Errorf("upgrade: %w", err)
	}
	defer conn.Close()

	tripID := c.Request.URL.Query().Get("trip_id")
	if tripID == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("No trip specified"))
		return fmt.Errorf("no trip specified")
	}

	// TODO: check if user is passenger for this trip

	_, err = locationws.ServeClient(
		conn,
		tripID,
		false,
	)
	if err != nil {
		return fmt.Errorf("serve client: %w", err)
	}

	return nil
}
