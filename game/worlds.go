package game

//World is an instance to run on the server
type World struct {
	Name    string
	Owner   Player
	Players []Player
}

//CreateWorld makes a slice of world instances
func CreateWorld() *World {
	return &World{}
}
