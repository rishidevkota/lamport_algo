package main

import (
	"fmt"
	"lamport/message"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const SERVER_PORT = 1235

type request struct {
	OwnerID   int
	TimeStamp int
}

var requestQue = make([]request, 0)
var clock int
var networkSize int
var nodeID int
var receivedReplies int

type Node int

func createRequest(duration int) {
	clock++
	requestQue = append(requestQue, request{
		OwnerID:   nodeID,
		TimeStamp: clock,
	})
	msg := message.Message{
		MessageType: message.REQUEST,
		SenderID:    nodeID,
		TimeStamp:   clock,
	}
	msg.Send(SERVER_PORT, true)
}

func simulateCriticalRegion() {
	fmt.Println("enter the critical region")
	time.Sleep(time.Millisecond * 8000)
	fmt.Println("exit the critical region")
}

func enterCriticalRegion() {
	receivedReplies = 0
	simulateCriticalRegion()
	msg := message.Message{
		MessageType: message.RELEASE,
		SenderID:    nodeID,
	}

	msg.Send(SERVER_PORT, true)
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
		msg.Send(req.OwnerID, false)
	} else if nodeHasPermissions() {
		//TODO
		requestQue = requestQue[1:]
		enterCriticalRegion()
	}
}

//adding on Request Que
func (n *Node) Request(msg *message.Message, reply *int) error {
	fmt.Printf("recive request from %d\n", msg.SenderID)
	requestQue = append(requestQue, request{
		OwnerID:   msg.SenderID,
		TimeStamp: msg.TimeStamp,
	})

	req := requestQue[0]
	if req.TimeStamp == msg.TimeStamp {
		_msg := message.Message{
			MessageType: message.REPLY,
			SenderID:    nodeID,
			TimeStamp:   msg.TimeStamp,
		}
		_msg.Send(msg.SenderID, false)
	}

	*reply = 0
	return nil
}

//confirming getting Critical Region
func (n *Node) Reply(msg *message.Message, reply *int) error {
	fmt.Printf("recive reply from %d\n", msg.SenderID)
	receivedReplies++
	if nodeHasPermissions() {
		fmt.Println("All permission were received")

		req := requestQue[0]
		if req.OwnerID == nodeID {
			requestQue = requestQue[1:]
			enterCriticalRegion()
		}
	}
	*reply = 0
	return nil
}

//removing item from Que
func (n *Node) Release(msg *message.Message, reply *int) error {
	fmt.Printf("recive release from %d\n", msg.SenderID)
	if len(requestQue) > 0 {
		requestQue = requestQue[1:]
	}

	processNextRequest()

	*reply = 0
	return nil
}

func (n *Node) NetworkSize(msg *message.Message, reply *int) error {
	fmt.Printf("recive network_size from %d\n", msg.SenderID)
	networkSize++

	*reply = 0
	return nil
}

func serve(l net.Listener, c chan string) {
	http.Serve(l, nil)
	c <- "done"
}

func main() {
	//First connect to server get id for this node/process
	//And ready to recive message from other node or server
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%d", SERVER_PORT))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var network message.Network
	err = client.Call("Server.RegisterNode", SERVER_PORT+1, &network)
	if err != nil {
		log.Fatal("server call error:", err)
	}

	nodeID = network.Port
	networkSize = network.Size

	node := new(Node)
	rpc.Register(node)
	rpc.HandleHTTP()

	log.Printf("Starting node at port: %d\n", nodeID)
	// err = http.ListenAndServe(fmt.Sprintf(":%d", reply), nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", nodeID))
	if e != nil {
		log.Fatal("listen error:", e)
	}

	c1 := make(chan string)
	c2 := make(chan int)
	go serve(l, c1)
	go func() {
		for {
			var input int
			fmt.Scanf("%d", &input)
			c2 <- input
		}
	}()
	for {
		select {
		case msg1 := <-c1:
			fmt.Println(msg1)
		case msg2 := <-c2:
			createRequest(msg2)
		}
	}
}
