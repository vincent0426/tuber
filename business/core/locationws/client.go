package locationws

import (
	"errors"
	"fmt"
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
		fmt.Println("sendLoop")
		select {
		case <-c.pingTicker.C:
			c.updateWriteDeadline()
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("failed to write ping message to the client %v", err)
				return
			}
		case msg := <-c.messageToSend:
			fmt.Println("receive msg from messageToSend")
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
		fmt.Println("receiveLoop")
		msgType, msg, err := c.ReadMessage()
		fmt.Println("msg:", string(msg))
		if err != nil {
			fmt.Println("err:", err)
			return // Usually caused by a normal client disconnection
		}

		fmt.Println("57")
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
	fmt.Println("doJoin")
	if c.broadcastRoom != nil {
		err := c.leaveRoom()
		if err != nil {
			return err
		}
	}

	timer := time.NewTimer(roomResponseTimeout)
	defer timer.Stop()

	broadcastRoom := AcquireBroadcastRoom(tripID)
	fmt.Println("AcquireBroadcastRoom Success")
	select {
	case broadcastRoom.channel.Joining <- c:
		c.broadcastRoom = broadcastRoom
		go c.listenRoomClose(broadcastRoom.channel.Closed)
		return nil
	case <-timer.C:
		return errors.New("join to broadcast room timeout")
	}
}
