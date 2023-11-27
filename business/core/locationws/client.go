package locationws

import (
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	*websocket.Conn
	pingTicker    *time.Ticker
	messageToSend chan string
	broadcastRoom *BroadcastRoom
}

// sendLoop is a loop to send messages to the websocket connection
func (c *Client) sendLoop() {
	defer c.Close()

	for {
		select {
		case <-c.pingTicker.C:
			c.updateWriteDeadline()
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("failed to write ping message to the client %v", err)
				return
			}
		case msg := <-c.messageToSend:
			c.updateWriteDeadline()
			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("failed to write message to the client %v", err)
				return
			}
		}
	}
}

func (c *Client) receiveLoop() {
	defer c.Close()

	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			return // Usually caused by a normal client disconnection
		}

		c.updateReadDeadline()

		if msgType == websocket.TextMessage { // only handle text message
			err := HandleRequest(c, msg)
			if err != nil {
				log.Printf("failed to handle request %v", err)
			}
			continue
		}
	}
}

func (c *Client) updateReadDeadline() {
	_ = c.SetReadDeadline(time.Now().Add(readMaxWait)) // ignore error
}

func (c *Client) updateWriteDeadline() {
	_ = c.SetWriteDeadline(time.Now().Add(writeMaxWait)) // ignore error
}

func newClient(baseConn *websocket.Conn) *Client {
	return &Client{
		Conn:          baseConn,
		pingTicker:    time.NewTicker(pingPeriod),
		messageToSend: make(chan string, msgChannelBufferSize),
	}
}

func (c *Client) listenRoomClose(closeChan chan struct{}) {
	select {
	case <-closeChan:
		c.Close()
		return
	}
}

func (c *Client) leaveRoom() error {
	broadcastRoom := c.broadcastRoom
	if broadcastRoom == nil {
		return nil
	}

	// If the room channel is already closed, then we don't need to send the
	// `Leaving` channel.
	if !broadcastRoom.channel.IsClosed() {
		timer := time.NewTimer(roomResponseTimeout)
		defer timer.Stop()

		select {
		case broadcastRoom.channel.Leaving <- c:
		case <-broadcastRoom.channel.Closed:
		case <-timer.C:
			return errors.New("leave from broadcast room timeout")
		}
	}

	c.broadcastRoom = nil
	return nil
}

func (c *Client) doJoin(
	tripID string,
) error {
	if c.broadcastRoom != nil {
		err := c.leaveRoom()
		if err != nil {
			return err
		}
	}

	timer := time.NewTimer(roomResponseTimeout)
	defer timer.Stop()

	broadcastRoom := AcquireBroadcastRoom(tripID)

	select {
	case broadcastRoom.channel.Joining <- c:
		c.broadcastRoom = broadcastRoom
		go c.listenRoomClose(broadcastRoom.channel.Closed)
		return nil
	case <-timer.C:
		return errors.New("join to broadcast room timeout")
	}
}
