package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sandblox-official/sockets-server/server"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var uid = 0

type worldName struct {
	Name string `json:"name"`
}

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
	resp, err := http.Get("http://localhost:8080/worlds")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var worldNames = []worldName{}
	json.Unmarshal(bodyBytes, &worldNames)
	for key, world := range worldNames {
		log.Println("Serve world", key+1, "/", len(worldNames), "[", world.Name, "]")
		worlds[world.Name] = server.CreateWorld()
		go func(world string) {
			worlds[world].Run()

		}(world.Name)
		world := world
		http.HandleFunc("/"+world.Name, func(w http.ResponseWriter, r *http.Request) {
			log.Println("Client connected to", world)
			serveWs(worlds[world.Name], w, r)
		})
	}

	//Serve and Run Worlds
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Println("")
	log.Println("Server ready on port", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(world *server.World, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
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
