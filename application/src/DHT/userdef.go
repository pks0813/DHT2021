package DHT

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
func FMSNewnode(port int)*DHTNode.DHTNode{
	Iter:=new(DHTNode.DHTNode)
	Iter.Data.On=true
	Iter.Data.SelfInfo.NexAddr= LocalAddress +":"+strconv.Itoa(port)
	Iter.Data.SelfInfo.Key=*DHTNode.PKSHash(Iter.Data.SelfInfo.NexAddr)
	tmp, err := net.Listen("tcp", Iter.Data.SelfInfo.NexAddr)
	fmt.Println("your IP is:"+Iter.Data.SelfInfo.NexAddr)
	if err != nil {
		fmt.Println("error in tools.go Run()")
		return nil
	}
	Iter.RPCServer=rpc.NewServer()
	Iter.Listen=new(net.Listener)
	err=Iter.RPCServer.Register(&(Iter.Data))
	if err!=nil{
		fmt.Println("error in userdef Newnode")
		return nil
	}
	Iter.Data.Data[0].MP=make(map[string]string)
	Iter.Data.Data[1].MP=make(map[string]string)
	Iter.Listen=&tmp
	go Iter.RPCServer.Accept(tmp)
	go Iter.Data.Update()
	return Iter

}
/*
func NewNode(port int) *DHTNode.DHTNode {
	Iter:=new(DHTNode.DHTNode)
	Iter.Data.SelfInfo.NexAddr= LocalAddress +":"+strconv.Itoa(port)
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
*/
// Todo: implement a struct which implements the interface "dhtNode".
