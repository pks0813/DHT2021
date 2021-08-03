package main

import (
	"DHTNode"
	"time"
)
var pksWait=100*time.Millisecond
func main() {
	//var pks1 DHTNode.DHTNode
	//pks1.Ping("f")
	//pks.Put("daya","1")
	var pks *DHTNode.DHTNode
	pks=NewNode(1)
	pks.Run()
	pks.Create()
	pks.Put("daya","shazi")
	pks.Put("wangzi","Big fool")
	pks.Put("pks","Smart!")
	pks.Put("wz","Big fool")
	pks.Put("pankaisen","Smart")
	time.Sleep(pksWait)
	var wz *DHTNode.DHTNode
	wz=NewNode(2)
	wz.Run()
	pks.Print()
	wz.Print()
	wz.Join(":1")
	time.Sleep(pksWait)
	//pks.Print()
	//wz.Print()
	wz.Quit()
	time.Sleep(pksWait*10)
	pks.Print()
}
