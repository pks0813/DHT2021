package main

import (
	"DHTNode"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)
/* In this file, you should implement function "NewNode" and
 * a struct which implements the interface "dhtNode".
 */

func NewNode(port int) dhtNode {
	Iter:=new(DHTNode.DHTNode)
	Iter.Data.SelfInfo.NexAddr=localAddress+":"+strconv.Itoa(port)
	Iter.Data.SelfInfo.Key=*DHTNode.PKSHash(Iter.Data.SelfInfo.NexAddr)
	Iter.RPCServer=rpc.NewServer()
	Iter.Listen=new(net.Listener)
	err:=Iter.RPCServer.Register(&(Iter.Data))
	Iter.Data.Data[0].MP=make(map[string]string)
	Iter.Data.Data[1].MP=make(map[string]string)
	if err!=nil{
		fmt.Println("error in userdef Newnode")
		panic(err)
	}
	return Iter
	// Todo: create a node and then return it.
}

// Todo: implement a struct which implements the interface "dhtNode".
