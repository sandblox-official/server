package server

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

//Worlds ...
var Worlds = make(map[string]*World)

//World ...
type World struct {
	Clients   map[*Client]bool
	Broadcast chan []byte
	Join      chan *Client
	Leave     chan *Client
}

//Client ...
type Client struct {
	ID    int
	World *World
	Conn  *websocket.Conn
	Send  chan []byte
}

//Emit allows for incoming messages
func (c *Client) Emit() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//Consume allows for outgoing messages
func (c *Client) Consume() {
	defer func() {
		c.World.Leave <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		messageNativePacket := &Packet{}
		err = json.Unmarshal([]byte(message), &messageNativePacket)
		if err != nil {
			log.Println("JSON Conversion err", err)
		}
		outMessage := messageNativePacket.GetOutputPacket()
		outJSONMessage, err := json.Marshal(outMessage)
		c.World.Broadcast <- outJSONMessage

	}
}

//CreateWorld ...
func CreateWorld() *World {
	return &World{
		Broadcast: make(chan []byte),
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Clients:   make(map[*Client]bool),
	}
}

//Run starts a continous loop for world updates
func (w *World) Run() {
	for {
		select {
		case client := <-w.Join:
			w.Clients[client] = true
		case client := <-w.Leave:
			if _, ok := w.Clients[client]; ok {
				delete(w.Clients, client)
				close(client.Send)
			}
		case message := <-w.Broadcast:
			for client := range w.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(w.Clients, client)
				}
			}
		}
	}
}
