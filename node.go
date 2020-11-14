package main

import (
	"fmt"
	"lamport/message"
	"time"
)

const SERVER_PORT = 1235

type request struct {
	OwnerID   int
	timeStamp int
}

var requestQue = make([]request, 0)
var clock int
var networkSize int
var nodeID int
var receivedReplies int

type Node int

func simulateCriticalRegion() {
	fmt.Println("enter the critical region")
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("exit the critical region")
}

func enterCriticalRegion() {
	receivedReplies = 0
	simulateCriticalRegion()
	msg := 
}

func nodeHasPermissions() bool {
	return receivedReplies == (networkSize - 1)
}

func processNextRequest() {
	if len(requestQue) == 0 {
		return
	}

	req := requestQue[0]
	if req.OwnerID != nodeID {
		msg := message.Message{
			MessageType: message.REPLY,
			SenderID:    nodeID,
			TimeStamp:   clock,
		}
		msg.Send(req.OwnerID)
	} else if nodeHasPermissions() {
		//TODO
		requestQue = requestQue[1:]
		fmt.Println("enter the critical region")
	}
}

//adding on Request Que
func (n *Node) Request(msg *message.Message, reply *int) error {
	//fmt.Println(msg.SenderID, "REQ")
	requestQue = append(requestQue, request{
		OwnerID:   msg.SenderID,
		timeStamp: msg.TimeStamp,
	})

	req := requestQue[0]
	if req.timeStamp == msg.TimeStamp {
		_msg := message.Message{
			MessageType: message.REPLY,
			SenderID:    nodeID,
			TimeStamp:   msg.TimeStamp,
		}
		_msg.Send(msg.SenderID)
	}

	*reply = 0
	return nil
}

//confirming getting Critical Region
func (n *Node) Reply(msg *message.Message, reply *int) error {
	//fmt.Println(msg.SenderID, "REP")
	receivedReplies++
	if nodeHasPermissions() {
		fmt.Println("All permission were received")

		req := requestQue[0]
		if req.OwnerID == nodeID {
			//TODO
			//enter critical region
			fmt.Println("enter critical region")
		}
	}
	*reply = 0
	return nil
}

//removing item from Que
func (n *Node) Release(msg *message.Message, reply *int) error {
	if len(requestQue) > 0 {
		requestQue = requestQue[1:]
	}

	processNextRequest()

	*reply = 0
	return nil
}

func main() {
	//First connect to server get id for this node/process
	//And ready to recive message from other node or server
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%d", SERVER_PORT))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply int
	err = client.Call("Server.RegisterNode", SERVER_PORT+1, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	node := new(Node)
	rpc.Register(node)
	rpc.HandleHTTP()

	log.Printf("Starting node at port: %d\n", reply)
	err = http.ListenAndServe(fmt.Sprintf(":%d", reply), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
