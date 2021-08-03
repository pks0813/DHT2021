package main

import (
	"DHT"
	"DHTNode"
	"FMS"
	"fmt"
	"os"
	"time"
)

func Init() {

}
const firstport=20000
func main(){
	os.Mkdir("Download",0666)
	os.Mkdir("Upload",0666)
	os.Mkdir("Seed",0666)
	fmt.Println("Welcome to File Management System")
	//_,_=fmt.Scanf("%d",&Num)
	fmt.Print("Which port do you want to use")
	var node *DHTNode.DHTNode
	var port int
	fmt.Scanf("%d",&port)
	node=DHT.FMSNewnode(port)
	for node==nil{
		fmt.Print("Given port are used, please give anohter port:")
		fmt.Scanf("%d",&port)
		node=DHT.FMSNewnode(port)
	}
	var op int
	fmt.Print("Node have build successfully do you want to Add a network(0) or Create a New network(1):")
	fmt.Scanf("%d",&op)
	if op==1{
		node.Create()
	}else {
		fmt.Print("Please given the IP+port of network:")
		var IP string
		fmt.Scanf("%s",&IP)
		node.Join(IP)
	}

	for true{
		time.Sleep(500*time.Millisecond)
		fmt.Println("What do you want to do(read your option)")
		fmt.Println("option 0: upload ")
		fmt.Println("option 1: download ")
		//fmt.Println("option 2: NewNode ")
		//fmt.Println("option 3: Quit ")
		fmt.Println("option 2: exit ")
		fmt.Scanf("%d",&op)
		switch op {
			case 0:FMS.Upload(node)
			case 1:FMS.Download(node)
			default:
				{
					node.Quit()
					time.Sleep(time.Second)
					fmt.Println("you have quit success")
					return
				}
		}
	}


}