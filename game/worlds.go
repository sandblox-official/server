package game

//World is an instance to run on the server
type World struct {
	Name      string
	Owner     Player
	Players   []Player
	Clients   map[*Client]bool
	Broadcast chan []byte
	Join      chan *Client
	Leave     chan *Client
}

//CreateWorld makes a slice of world instances
func CreateWorld(name string, owner Player) *World {
	return &World{Name: name, Owner: owner}
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
