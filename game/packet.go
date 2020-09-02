package game

//Packet ...
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
		outPacket.Data.Player = inPacket.Data.Player
		return *outPacket
	case "message":
		outPacket.Data.Chat = inPacket.Data.Chat
		return *outPacket
	}

	outPacket.Method = "error"
	return *outPacket
}
