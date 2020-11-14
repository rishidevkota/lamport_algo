package message

import (
	"fmt"
	"log"
	"net/rpc"
)

const (
	REQUEST = iota
	REPLY
	RELEASE
	NETWORK_SIZE
)

type Message struct {
	MessageType int
	SenderID    int
	TimeStamp   int
}

//Broadcast message or Send message to particular process
func (msg *Message) Send(to int) interface{} {
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%d", to))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// msg := Message{
	// 	SenderID: from,
	// }

	// var reply int
	// err = client.Call("Node.Reply", msg, &reply)
	// if err != nil {
	// 	log.Fatal("arith error:", err)
	// }

	// fmt.Println(reply)
	//var method string
	var reply int
	switch msg.MessageType {
	case REQUEST:
		err = client.Call("Node.Request", msg, &reply)
		if err != nil {
			log.Fatal("node request send error:", err)
		}
	case REPLY:

	case RELEASE:

	case NETWORK_SIZE:

	}

	return reply
}
