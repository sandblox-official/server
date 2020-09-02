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
	//Static Files
	fs := http.FileServer(http.Dir("./webroot"))
	http.Handle("/", fs)
	//Socket Hanlder
	worlds := server.Worlds
	worlds["test1"] = server.CreateWorld()
	go worlds["test1"].Run()
	http.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Client connected to world 1")
		serveWs(worlds["test1"], w, r)
	})
	worlds["test2"] = server.CreateWorld()
	go worlds["test2"].Run()
	http.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Client connected to world 2")
		serveWs(worlds["test2"], w, r)
	})
	//Serve and Run Worlds
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
