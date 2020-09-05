package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sandblox-official/game/server"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var uid = 0

func main() {
	//Logs
	f, err := os.OpenFile("./logs/main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	//Get all worlds from database
	worlds := server.Worlds
	worldsFromDB := []string{
		"world1",
		"world2",
		"world3",
	}
	for _, world := range worldsFromDB {
		worlds[world] = server.CreateWorld()
		go func(world string) {
			worlds[world].Run()

		}(world)
		world := world
		http.HandleFunc("/"+world, func(w http.ResponseWriter, r *http.Request) {
			log.Println("Client connected to", world)
			serveWs(worlds[world], w, r)
		})
	}

	//Serve and Run Worlds
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(world *server.World, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &server.Client{ID: uid, World: world, Conn: conn, Send: make(chan []byte, 256)}
	client.World.Join <- client
	uid++

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.Emit()
	go client.Consume()

}
