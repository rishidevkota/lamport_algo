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
func (msg *Message) Send(to int, isBroadcast bool) interface{} {
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
		if isBroadcast {
			err = client.Call("Server.Request", msg, &reply)
		} else {
			err = client.Call("Node.Request", msg, &reply)
		}
		if err != nil {
			log.Fatal("request send error:", err)
		}
	case REPLY:
		err = client.Call("Node.Reply", msg, &reply)
		if err != nil {
			log.Fatal("reply send error:", err)
		}
	case RELEASE:
		if isBroadcast {
			err = client.Call("Server.Release", msg, &reply)
		} else {
			err = client.Call("Node.Release", msg, &reply)
		}
		if err != nil {
			log.Fatal("release send error:", err)
		}
	case NETWORK_SIZE:
		err = client.Call("Node.NetworkSize", msg, &reply)
		if err != nil {
			log.Fatal("network_size send error:", err)
		}
	}

	return reply
}
