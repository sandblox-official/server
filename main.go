package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sandblox-official/server/game"
)

var upgrader = websocket.Upgrader{}
var uid = 1

func main() {
	//Set log file
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Server start")
	worlds := make(map[string]*game.World)
	player1 := game.Player{
		Name: "player1",
		X:    0.5,
		Y:    3,
		Z:    1,
	}
	worlds["test1"] = game.CreateWorld("test1", player1)
	go worlds["test1"].Run()
	http.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Client connected to world 1")
		serveWs(worlds["test1"], w, r)
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
func serveWs(world *game.World, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &game.Client{ID: uid, World: world, Conn: conn, Send: make(chan []byte, 256)}
	client.World.Join <- client
	uid++

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.Emit()
	go client.Consume()

}
