package game

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

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
