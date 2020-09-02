package game

import "github.com/gorilla/websocket"

//Player ..
type Player struct {
	Name string
	X    float32
	Y    int
	Z    int
}

//Client is the player instance coming in
type Client struct {
	ID    int
	World *World
	Conn  *websocket.Conn
	Send  chan []byte
}
