package server

//Packet is the data format for socket broadcasts
type Packet struct {
	Method string `json:"method"`
	Data   Data
}

//Data ...
type Data struct {
	//Outgoing
	Player Player
	Chat   Chat
}

//Player ...
type Player struct {
	Name string
	X    float32
	Y    int
	Z    int
}

//Chat ...
type Chat struct {
	From string
	Body string
}

//GetOutputPacket takes an input to generate an output
func (inPacket *Packet) GetOutputPacket() Packet {
	outPacket := &Packet{}
	switch inPacket.Method {
	case "move":
		outPacket.Method = "move"
		outPacket.Data.Player = inPacket.Data.Player
		return *outPacket
	case "message":
		outPacket.Data.Chat = inPacket.Data.Chat
		return *outPacket
	}

	outPacket.Method = "error"
	return *outPacket
}
