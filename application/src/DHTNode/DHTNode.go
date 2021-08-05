package DHTNode

import (
	"fmt"
	"math/big"
	"net"
	"net/rpc"
	"sync"
	"time"
)

var FMSflag bool
type Replytype struct {
	Begin,End int
	Ans []byte
}
type Asktype struct {
	Begin,End int
	Ask []byte
}
type DHTNode struct {
	Data Node
	Listen *net.Listener
	RPCServer *rpc.Server
}
func (this *DHTNode)FMSGet(Key Asktype,reply *chan Replytype){
	var tmp Replytype
	tmp.Begin=Key.Begin
	tmp.End=Key.End
	pks1,Back:=this.Get(string(Key.Ask))
	if pks1==false{
		FMSflag=false
		fmt.Println("Download Fail")
		*reply<-tmp
	}
	tmp.Ans=[]byte(Back)
	*reply<-tmp
}
func (this *DHTNode)FMSPut(Hash string,value []byte,sum *int,zong *int,wg *sync.WaitGroup,time0 *time.Time){
	//fmt.Println("Put:",Hash,value)
	defer func(){
		*sum+=len(value)
		fmt.Println("Loading",float32(*sum)/1024/1024,"MB/",float32(*zong)/1024/1024,"MB        percent:",float32(*sum)/float32(*zong)*100,"%     NowTime:",time.Since(*time0))
		wg.Done()
	}()
	this.Put(Hash,string(value))

}
func (this *DHTNode)Run(){
	tmp, err := net.Listen("tcp", this.Data.SelfInfo.NexAddr)
	if err != nil {
		fmt.Println("error in tools.go Run()")
 		panic(nil)
	}
	this.Data.On=true
	this.Listen=&tmp
	go this.RPCServer.Accept(tmp)
	go this.Data.Update()
}
func (this *DHTNode)Create(){
	this.Data.Lock.Lock()
	this.Data.PreInfo.NexAddr=this.Data.SelfInfo.NexAddr
	this.Data.PreInfo.Key.Set(&(this.Data.SelfInfo.Key))
	for i:=0;i<NexLen;i++{
		this.Data.NexInfo[i].NexAddr=this.Data.SelfInfo.NexAddr
		this.Data.NexInfo[i].Key.Set(&(this.Data.SelfInfo.Key))
	}
	for i:=0;i<M;i++{
		this.Data.FingerTable[i].NexAddr=this.Data.SelfInfo.NexAddr
		this.Data.FingerTable[i].Key.Set(&(this.Data.SelfInfo.Key))
	}
	this.Data.Circle=true
	this.Data.Lock.Unlock()
}

func (this *DHTNode)Join(Addr string) bool{
	Client,err:=Dial(Addr)
	
	if err!=nil{
		fmt.Println(err)
		fmt.Println("error in join Can't find given Addr")
		return false
	}
	for i:=0;i<M;i++{
		var Nowpos big.Int
		Nowpos.Add(&this.Data.SelfInfo.Key,new(big.Int).Exp(big.NewInt(2),big.NewInt(int64(i)),nil))
		Nowpos.Mod(&Nowpos,Mod)
		var tmp EdgePair
		err=Client.Call("Node.FindSuccessor",&Nowpos,&tmp)
		this.Data.FingerTable[i]=tmp
		if err!=nil{
			fmt.Println(err)
			fmt.Println("error in join to Fixfinger",i)
		}
	}
	var tmp EdgePair
	err=Client.Call("Node.FindSuccessor",&this.Data.SelfInfo.Key,&tmp)
	_=Client.Close()
	DialTime--
	if err!=nil{
		fmt.Println(err)
		fmt.Println("error in join Can't find Successor")
		return false
	}
	this.Data.insSuc(&tmp)
	Client,err=Dial(tmp.NexAddr)
	if err!=nil{
		fmt.Println(err)
		return false
	}
	var Nowmp map[string]string
	err=Client.Call("Node.Split",&this.Data.SelfInfo.Key,&Nowmp)
	this.Data.Data[0].MP=Nowmp

	//err=Client.Call("Node.Notify",this.Data.SelfInfo,nil)

	_=Client.Close()
	DialTime--
	this.Data.Circle=true
	return true
}
func (this *DHTNode)Quit(){
	this.Data.On=false
	this.Data.Circle=false
	_=(*this.Listen).Close()
	//fmt.Println(err)
	//fmt.Println("Quit Finish")
}
func (this *DHTNode)ForceQuit(){
	this.Data.On=false
	this.Data.Circle=false
	_=(*this.Listen).Close()
}
func (this *DHTNode)Ping(Addr string)bool{
	return this.Data.ping(Addr)
}

func (this *DHTNode) Put(key string, value string) bool{
	Hash:=PKSHash(key)
	var tmp EdgePair
	err:=this.Data.FindSuccessor(Hash,&tmp)
	//return true
	if err!=nil{
		fmt.Println("error in put Cant't Find Client")
		return false
	}
	Client,err:=Dial(tmp.NexAddr)
	if err!=nil{
		fmt.Println("error in put Dial Fail",tmp.NexAddr)
		return false
	}
	err=Client.Call("Node.InsertKV",KVPair{key,value},nil)
	_=Client.Close()
	DialTime--
	if err!=nil{
		fmt.Println("can't Call Success")
		return false
	}
	return true
}
func (this *DHTNode)Get(key string) (bool, string){
	Hash:=PKSHash(key)
	var tmp EdgePair
	err:=this.Data.FindSuccessor(Hash,&tmp)
	if err!=nil{
		fmt.Println("error in put Cant't Find Client")
		return false,""
	}
	var value string
	Client,err:=Dial(tmp.NexAddr)
	
	if err!=nil{
		fmt.Println("can't Call Success")
		return false,""
	}
	err=Client.Call("Node.Query",key,&value)
	_=Client.Close()
	DialTime--
	if err!=nil{
		return false,""
	}
	return true,value
}
func (this *DHTNode) Delete(key string) bool{
	Hash:=PKSHash(key)
	var tmp EdgePair
	err:=this.Data.FindSuccessor(Hash,&tmp)
	if err!=nil{
		fmt.Println("error in put Cant't Find Client")
		return false
	}
	Client,err:=Dial(tmp.NexAddr)

	if err!=nil{
		return false
	}
	err=Client.Call("Node.Delete",key,nil)
	_=Client.Close()
	DialTime--
	return true
}

func (this *DHTNode)Print(){
	fmt.Println("SelfInfomation:",&this.Data.SelfInfo.Key,this.Data.SelfInfo.NexAddr)
	fmt.Println("PreInfomation:",&this.Data.PreInfo.Key,this.Data.PreInfo.NexAddr)
	fmt.Println("NexInfomation:",&this.Data.NexInfo[0].Key,this.Data.NexInfo[0].NexAddr)
	//fmt.Println("Map:")
	//for key,value:=range this.Data.Data[0].MP{fmt.Println(PKSHash(key),"    Key:",key,"	Value:",value)}
	//fmt.Println("PreMap:")
	//for key,value:=range this.Data.Data[1].MP{fmt.Println(PKSHash(key),"    Key:",key,"	Value:",value)}
	//for i:=0;i<M;i++{
	//	fmt.Println("Fingertable :",i,this.Data.FingerTable[i].NexAddr,this.Data.FingerTable[i].Key)
}
