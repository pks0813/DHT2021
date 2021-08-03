package DHTNode

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"
)
const M int=10
const NexLen int=20
const UpDateTime = 100 * time.Millisecond
type EdgePair struct{
	Key big.Int
	NexAddr string
}

type Storage struct {
	MP map[string]string
	BooLock sync.Mutex
}

type Node struct{
	NexInfo [NexLen]EdgePair
	PreInfo EdgePair
	SelfInfo EdgePair
	Lock sync.Mutex
	On bool
	FingerTable [M]EdgePair
	Circle bool
	Data[2] Storage
	fix int
}

const WaitTime= time.Millisecond

func (this * Node) Uppredata(tmp *map[string]string,_ *int)error{
	this.Data[1].BooLock.Lock()
	this.Data[1].MP=*tmp
	this.Data[1].BooLock.Unlock()
	return nil
}

func (this *Node) CheckPre(_ int,_ *int) error{
	//fmt.Println("checkPre",this.SelfInfo.NexAddr)
	this.Lock.Lock()
	PreIP:=this.PreInfo.NexAddr
	this.Lock.Unlock()
	if PreIP!="" && this.ping(PreIP)==false {
		this.Data[0].BooLock.Lock()
		this.Data[1].BooLock.Lock()
		PKSMerge(&(this.Data[0].MP), &(this.Data[1].MP))
		this.Data[1].MP=make(map[string]string)
		this.Data[0].BooLock.Unlock()
		this.Data[1].BooLock.Unlock()
		if this.fixList() != nil {
			fmt.Println("Error in CheckPre in FixList")
			return errors.New("Error in CheckPre in FixList")
		}
		this.Lock.Lock()
		NexIP := this.NexInfo[0].NexAddr
		this.Lock.Unlock()
		NexClient, err := Dial(NexIP)
		if err != nil {
			fmt.Println("error in Checkpre() ", err)
			return nil
		}
		this.Data[0].BooLock.Lock()
		tmp := this.Data[0].MP
		this.Data[0].BooLock.Unlock()
		err = NexClient.Call("Node.Uppredata", &tmp, nil)
		_ = NexClient.Close()
		DialTime--
		if err != nil {
			fmt.Println("error in Checkpre() ", err)
			return nil
		}
		this.Lock.Lock()
		this.PreInfo.NexAddr = ""
		this.Lock.Unlock()
	}
	return nil
}
func (this *Node)DataCopy(_ int,Info *map[string]string)error{//把自己的storage复制出来
	this.Data[0].BooLock.Lock()
	*Info=this.Data[0].MP
	this.Data[0].BooLock.Unlock()
	return nil
}
func (this *Node)Notify(Info *EdgePair,_ *int)error{
	//fmt.Println("Notify",Info.Key,Info.NexAddr)
	this.Lock.Lock()
	Pre:=this.PreInfo
	Self:=this.SelfInfo
	this.Lock.Unlock()
	if Pre.NexAddr==""{
		this.Lock.Lock()
		this.PreInfo=*Info
		this.Lock.Unlock()
		Client,err:=Dial(Info.NexAddr)
		if err!=nil{
			fmt.Println(err,Info.NexAddr)
			return errors.New("error in Notify")
		}
		var tmp map[string]string
		err=Client.Call("Node.DataCopy",0,&tmp)
		_=Client.Close()
		DialTime--
		this.Data[1].BooLock.Lock()
		this.Data[1].MP = tmp
		this.Data[1].BooLock.Unlock()
	}else {
		if Inrange(&(Pre.Key),&(Self.Key),&(Info.Key)){
			this.Lock.Lock()
			this.PreInfo=*Info
			this.Lock.Unlock()
			Client,err:=Dial(Info.NexAddr)
			if err!=nil{
				fmt.Println(err,Info.NexAddr)
				return errors.New("error in Notify")
			}
			var tmp map[string]string
			err=Client.Call("Node.DataCopy",0,&tmp)
			_=Client.Close()
			DialTime--
			this.Data[1].BooLock.Lock()
			this.Data[1].MP = tmp
			this.Data[1].BooLock.Unlock()
		}
	}
	return nil
}
func (this *Node) GetPre(_ int ,reply *EdgePair)error {
	this.Lock.Lock()
	tmp:=this.PreInfo
	this.Lock.Unlock()
	*reply=tmp
	return nil
}
func (this *Node) changeDataPre(tmp *map[string]string,_ *int){
	this.Data[1].BooLock.Lock()
	this.Data[1].MP=*tmp
	this.Data[1].BooLock.Unlock()
}
func (this *Node) insSuc(tmp *EdgePair){
	this.Lock.Lock()
	for i:=NexLen-1;i>=1;i--{
		this.NexInfo[i]=this.NexInfo[i-1]
	}
	this.NexInfo[0]=*tmp
	this.Lock.Unlock()
	//update MP
	//好像因为后面会notify不用更新
	//Client,err:=Dial(tmp.NexAddr)
	//if err!=nil{
	//	fmt.Println("inssuc fail")
	//	return
	//}
	//var Q Storage
	//this.DataCopy(0,&Q)
	//Client.Call("Node.changeDataPre",&Q,nil)
	//Client.Close()
}

func (this *Node) CopyList(_ int,reply *[NexLen]EdgePair)error{
	this.Lock.Lock()
	for i:=0;i<NexLen;i++{
		(*reply)[i]=this.NexInfo[i]
	}
	this.Lock.Unlock()
	return nil
}
func (this *Node) Stabilize(_ int,_ *int)error {
	//fmt.Println("stabilize",this.SelfInfo.NexAddr)
	if this.fixList()!=nil{
		fmt.Println("Stabilize All suc fail")
		return nil
	}
	//fmt.Println("Fix Finish")
	this.Lock.Lock()
	NexInfo:=this.NexInfo[0]
	this.Lock.Unlock()
	//fmt.Println("nexInfo is",NexInfo.NexAddr)
	Client,err:=Dial(NexInfo.NexAddr)
	
	if err!=nil{
		fmt.Println(err)
		return nil
	}
	var tmp EdgePair
	err=Client.Call("Node.GetPre",0,&tmp)
	//fmt.Println("tmp is ",tmp.NexAddr)
	_=Client.Close()
	DialTime--
	//fmt.Println("GetNexPre Finish")
	if err!=nil{
		fmt.Println("error in Stabilize Getpre",err)
		return nil
	}
	this.Lock.Lock()
	SelfInfo:=this.SelfInfo
	this.Lock.Unlock()
	if Inrange(&(SelfInfo.Key),&(NexInfo.Key),&(tmp.Key))&& this.ping(tmp.NexAddr) {
		this.insSuc(&tmp)
	}
	this.Lock.Lock()
	NexInfo=this.NexInfo[0]
	this.Lock.Unlock()
	var NexList [NexLen]EdgePair
	Client,err=Dial(NexInfo.NexAddr)
	
	if err!=nil{
		fmt.Println("error in Stabilize ",err)
		fmt.Println(tmp.NexAddr)
		return nil
	}
	err=Client.Call("Node.CopyList",0,&NexList)
	this.Lock.Lock()
	for i:=1;i<NexLen-1;i++{
		this.NexInfo[i]=NexList[i-1]
	}
	this.Lock.Unlock()
	err=Client.Call("Node.Notify",&SelfInfo,nil)
	if err!=nil{
		fmt.Println("error in notify")
	}
	_=Client.Close()
	DialTime--
	return nil
}

func (this *Node)Query(Key string,reply *string)error{
	//var value string
	//var ok bool
	this.Data[0].BooLock.Lock()
	value,ok:=this.Data[0].MP[Key]
	this.Data[0].BooLock.Unlock()
	if ok{
		*reply=value
	}else{
		return errors.New("Query Fail")
	}
	return nil
}
func (this *Node)InsertKV(Ins KVPair,_ *int)error{
	this.Data[0].BooLock.Lock()
	this.Data[0].MP[Ins.Key]=Ins.Value
	this.Data[0].BooLock.Unlock()
	if this.fixList()!=nil{
		fmt.Println("error in InsertKV Can't Find any Succesor")
		return errors.New("error in InsertKV Can't Find any Succesor")
	}
	this.Lock.Lock()
	Client,err:=Dial(this.NexInfo[0].NexAddr)

	this.Lock.Unlock()
	if err!=nil{
		fmt.Println("error in InsertKV Can't Dial Succesor[0]")
		return errors.New("error in InsertKV Can't Dial Succesor[0]")
	}
	err=Client.Call("Node.InsertKVpre",Ins,nil)
	_=Client.Close()
	DialTime--
	return nil
}
func (this *Node)InsertKVpre(Ins KVPair,_ *int)error{
	this.Data[1].BooLock.Lock()
	this.Data[1].MP[Ins.Key]=Ins.Value
	this.Data[1].BooLock.Unlock()
	return nil
}
func (this *Node)Delete(Del string,_ *int)error{
	this.Data[0].BooLock.Lock()
	_,ok:=this.Data[0].MP[Del]
	if !ok{
		this.Data[0].BooLock.Unlock()
		return errors.New("Can't find Delete key")
	}
	delete(this.Data[0].MP,Del)
	this.Data[0].BooLock.Unlock()
	if this.fixList()!=nil{
		fmt.Println("error in InsertKV Can't Find any Succesor")
		return errors.New("error in InsertKV Can't Find any Succesor")
	}
	this.Lock.Lock()
	Client,err:=Dial(this.NexInfo[0].NexAddr)

	this.Lock.Unlock()
	if err!=nil{
		fmt.Println("error in InsertKV Can't Dial Succesor[0]")
		return errors.New("error in InsertKV Can't Dial Succesor[0]")
	}
	err=Client.Call("Node.Deletepre",Del,nil)
	_=Client.Close()
	DialTime--
	return nil
}
func (this *Node)Deletepre(Del string,_ *int)error{
	this.Data[0].BooLock.Lock()
	_,ok:=this.Data[0].MP[Del]
	if !ok{
		this.Data[0].BooLock.Unlock()
		return errors.New("Can't find Delete key")
	}
	delete(this.Data[0].MP,Del)
	this.Data[0].BooLock.Unlock()
	return nil
}
func (this *Node)fixList() error{
	//fmt.Println("FixList:",this.SelfInfo.NexAddr,this.NexInfo[0].NexAddr)
	this.Lock.Lock()
	p:=-1
	for i:=0;i<NexLen;i++{
		//fmt.Println("I Have Ping",this.NexInfo[i].NexAddr)
		if this.ping(this.NexInfo[i].NexAddr)==true{
			p=i
			break
		}
		//fmt.Println("FixListFind Fail :",this.SelfInfo.NexAddr,this.NexInfo[0].NexAddr)
	}
	//fmt.Println("FixList Ping Finish")
	if p==-1{
		this.Lock.Unlock()
		return errors.New("error in fixList All suc Fail")
	}
	if p==0{
		this.Lock.Unlock()
		return nil
	}
	for i:=0;i<NexLen-p;i++{
		this.NexInfo[i]=this.NexInfo[i+p]
	}
	//tmp:=this.NexInfo[0]
	//NowEdge:=this.SelfInfo
	this.Lock.Unlock()
	/*
	Client,err:=Dial(tmp.NexAddr)
	if err!=nil{
		fmt.Println("error in fixList",err,"Addr is ",tmp.NexAddr)
		return err
	}
	err=Client.Call("Node.Notify",&NowEdge,nil)
	_=Client.Close()
	DialTime--*/
	return nil
}
func (this *Node)Split(Prekey *big.Int,reply *map[string]string)error{
	this.Data[0].BooLock.Lock()
	this.Data[1].BooLock.Lock()
	*reply=make(map[string]string)
	for key,value:=range this.Data[0].MP{
		if !Inrange(Prekey,&this.SelfInfo.Key,PKSHash(key)){
			(*reply)[key]=value
		}
	}
	this.Data[1].MP=*reply
	for key,_:=range *reply{
		delete(this.Data[0].MP,key)
	}
	this.Data[0].BooLock.Unlock()
	this.Data[1].BooLock.Unlock()
	return nil
}
func (this *Node)FixFinger(_ int,_ *int)error{
	//fmt.Println("fixFinger",this.SelfInfo.NexAddr)
	this.fix=(this.fix+1)%M
	Fix:=this.fix
	var tmp big.Int
	var UpFix EdgePair
	tmp.Add(&(this.SelfInfo.Key),new(big.Int).Exp(big.NewInt(2),big.NewInt((int64)(Fix)),nil))
	tmp.Mod(&tmp,Mod)
	err:=this.FindSuccessor(&tmp,&UpFix)
	if err!=nil{
		fmt.Println("fixFinger Fail",err)
		return nil
	}
	this.Lock.Lock()
	this.FingerTable[Fix]=UpFix
	this.Lock.Unlock()
	return nil
}
func (this *Node) FindSuccessor(Addr *big.Int,reply *EdgePair)error{
	Del:=make(map[string]bool)
	if this.SelfInfo.Key.Cmp(Addr)==0{
		*reply=this.SelfInfo
		return nil
	}
	//fmt.Println("Findsuccessor:",this.SelfInfo.Key,this.NexInfo[0].Key,Addr)
	if this.fixList()!=nil{
		return errors.New("Find Successor Fail in FixList")
	}
	this.Lock.Lock()
	NexEdge:=this.NexInfo[0]
	this.Lock.Unlock()
	if this.ping(NexEdge.NexAddr)==false{
		fmt.Println("FindSuccessor Fail in Pingnex",this.SelfInfo.NexAddr,NexEdge.NexAddr)
		return errors.New("FindSuccesor Fail")
	}
	if Inrange(&(this.SelfInfo.Key),&(NexEdge.Key),Addr){
		*reply=NexEdge
		return nil
	}
	tmp:=NexEdge
	this.Lock.Lock()
	for i:=M-1;i>=0;i--{
		_,pks2:=Del[this.FingerTable[i].NexAddr]
		if pks2{continue}
		if !this.ping(this.FingerTable[i].NexAddr){
			Del[this.FingerTable[i].NexAddr]=true
			continue
		}
		if Inrange(&(this.SelfInfo.Key),Addr,&(this.FingerTable[i].Key)) {
			tmp=this.FingerTable[i]
			break
		}
	}
	this.Lock.Unlock()
	Client,err:=Dial(tmp.NexAddr)
	
	if tmp.NexAddr==this.SelfInfo.NexAddr{
		panic("you are fw")
	}
	if err!=nil{
		fmt.Println("error in FindSuccessor in Can't Dial nexIP")
		fmt.Println(err,tmp.NexAddr)
		return errors.New("error in FindSuccessor")
	}
	err=Client.Call("Node.FindSuccessor",Addr,reply)
	_=Client.Close()
	DialTime--
	if err!=nil{
		fmt.Println("FindSuccessor Fail")
		return errors.New("FindSuccesor Fail")
	}
	return nil
}
func (this *Node) Update() {
	for this.On{
		if this.Circle {
			this.CheckPre(0,nil)
			this.Stabilize(0,nil)
			this.FixFinger(0,nil)
		}
		time.Sleep(UpDateTime)
	}
}
func (this *Node)Update1(_ int ,_ *int)error{
	this.CheckPre(0,nil)
	this.Stabilize(0,nil)
	this.FixFinger(0,nil)
	return nil
}
