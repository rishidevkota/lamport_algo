//Purpose for this server is
//Creating id of node
//Broadcasting message to all node
//Updating size of network
package main

import (
	"fmt"
	"lamport/message"
	"net/http"
	"net/rpc"
)

const PORT = 1235

var NodePorts = make([]int, 0)

type Server int

// func (s *Server) Broadcast() error {
// 	return nil
// }

func broadcast(msg message.Message) {
	for _, _port := range NodePorts {
		msg.Send(_port)
	}
}

func (s *Server) Request(msg *message.Message, reply *int) error {

}

func (s *Server) Reply(msg *message.Message, reply *int) error {

}

func (s *Server) RegisterNode(port *int, reply *int) error {
	var isPortUsed bool
	for _, _port := range NodePorts {
		if _port == *port {
			isPortUsed = true
		}
	}
	if isPortUsed {
		*reply = NodePorts[len(NodePorts)-1] + 1
	} else {
		*reply = *port
	}
	msg := message.Message{
		MessageType: message.REQUEST,
		SenderID:    PORT,
	}
	broadcast(msg)

	NodePorts = append(NodePorts, *reply)

	return nil
}

// func (s *Server) GetNetworkSize() error {
// 	return nil
// }

func main() {
	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
