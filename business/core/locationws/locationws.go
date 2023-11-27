package locationws

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	//Todo: Mongo DB
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

const (
	// The upper bound on the total number of senders:
	// 1. the response for the websocket request
	// 2. the broadcast from the room
	msgChannelBufferSize = 2

	// Time allowed to write a message to the peer.
	writeMaxWait = 20 * time.Second

	// Send pings to peer with this period.
	pingPeriod  = 20 * time.Second
	pingTimeout = 20 * time.Second

	// Time allowed to read the next message from the peer.
	readMaxWait = pingPeriod + pingTimeout

	roomResponseTimeout = 2 * time.Minute
)

// ServeClient creates a new websocket connection and handles the corresponding interactions.
func ServeClient(
	baseConn *websocket.Conn,
	tripID string,
	isDriver bool,
) (*Client, error) {
	client := newClient(baseConn)

	client.SetCloseHandler(func(int, string) error {
		return nil
	})

	client.updateReadDeadline()
	client.SetPongHandler(func(string) error {
		client.updateReadDeadline()
		return nil
	})

	err := client.doJoin(tripID)
	if err != nil {
		return nil, err
	}

	go client.sendLoop()

	if isDriver {
		go client.receiveLoop()
	}

	return client, nil
}

func SendMessage(c *Client, msg string) error {
	timer := time.NewTimer(writeMaxWait)
	defer timer.Stop()

	select {
	case c.messageToSend <- msg:
		return nil
	case <-timer.C:
		return errors.New("failed to write into message-to-send channel on time")
	}
}

func HandleRequest(c *Client, plainReq []byte) error {
	var req Location
	err := json.Unmarshal(plainReq, &req)
	if err != nil {
		return err
	}

	ctx := context.Background()
	msg, err := json.Marshal(req)
	if err != nil {
		return err
	}

	roomsDispatcher.broadcastRoomMap[c.broadcastRoom.id].PublishMessage(ctx, string(msg))

	return nil
}
