package server

import "log"

//Packet is the data format for socket broadcasts
type Packet struct {
	Method string `json:"method"`
	Data   Data   `json:"data"`
}

//Data ...
type Data struct {
	//Outgoing
	Player Player `json:"player"`
	Chat   Chat   `json:"chat"`
}

//Player ...
type Player struct {
	Name string `json:"name"`
	X    float32
	Y    int
	Z    int
}

//Chat ...
type Chat struct {
	From string `json:"from"`
	Body string `json:"body"`
}

//GetOutputPacket takes an input to generate an output
func (inPacket *Packet) GetOutputPacket() Packet {
	outPacket := &Packet{}
	switch inPacket.Method {
	case "move":
		outPacket.Method = "move"
		outPacket = inPacket
		return *outPacket
	case "message":
		outPacket = inPacket
		log.Println("New message: '", inPacket.Data.Chat.Body, "'", "from", inPacket.Data.Chat.From)
		return *outPacket
	}

	outPacket.Method = "error"
	return *outPacket
}
