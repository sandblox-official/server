package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/sandblox-official/server/game"
)

func main() {
	log.Println("Server start")
	worlds := make(map[string]*game.World)
	spew.Dump(worlds)
}
