package FMS

import (
	"DHTNode"
	"bencode"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)
type FileInfo struct{
	Name string
	PieceHash string
	Length int
}

const PieceSize=1000000
func Download(Node *DHTNode.DHTNode){
	fmt.Print("What do you want to Download(link file name):")
	var inFile string
	fmt.Scanf("%s",&inFile)
	file,err:=os.Open("Seed/"+inFile)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//fmt.Println(string(data))
	var BT FileInfo
	err=bencode.Unmarshal(file,&BT)
	file.Close()
	if err!=nil{
		fmt.Println(err)
		return
	}
	BTByte:=[]byte(BT.PieceHash)
	//fmt.Println(BTByte)
	reply:=make([]byte,BT.Length)
	result:=make(chan DHTNode.Replytype)
	j:=0
	for i:=0;i<BT.Length;i+=PieceSize{
		var tmp DHTNode.Asktype
		tmp.Begin=i
		tmp.End=DHTNode.Pksmin(i+PieceSize,BT.Length)
		tmp.Ask=BTByte[j:j+20]
		j+=20
		go Node.FMSGet(tmp,&result)
	}
	sum:=0
	time0:=time.Now()
	for i:=0;i<BT.Length;i+=PieceSize{
		now:=<-result
		//time.Sleep(time.Millisecond*10)
		sum+=now.End-now.Begin
		fmt.Println("Loading",float32(sum)/1024/1024,"MB/",float32(BT.Length)/1024/1024,"MB        percent:",float32(sum)/float32(BT.Length)*100,"%     NowTime:",time.Since(time0))
		//fmt.Println("Loading",sum,"/",BT.Length,"           percent:",float32(sum)/float32(BT.Length)*100,"%")
		copy(reply[now.Begin:now.End],now.Ans)
	}
	file,err=os.Create("Download/"+BT.Name)
	file.Write(reply)
}
func Upload(Node *DHTNode.DHTNode){
	fmt.Print("What do you want to Upload:")
	var Reply FileInfo
	fmt.Scanf("%s",&Reply.Name)
	data,err := ioutil.ReadFile("Upload/"+Reply.Name)
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Print("Where do you want to store the BTseed:")
	var output string
	fmt.Scanf("%s",&output)
	Reply.Length=len(data)
	var wg sync.WaitGroup
	var sum int
	time0:=time.Now()
	for i:=0;i<Reply.Length;i+=PieceSize {
		low := i
		high := DHTNode.Pksmin(Reply.Length, i+PieceSize)
		h := sha1.New()
		h.Write(data[low:high])
		Hash := string(h.Sum(nil))
		wg.Add(1)
		go Node.FMSPut(Hash, data[low:high],&sum,&Reply.Length,&wg,&time0)
		Reply.PieceHash += Hash
	}
	wg.Wait()
	file,err:=os.Create("Seed/"+output)
	err=bencode.Marshal(file,Reply)
	//fmt.Println(Reply.Length,Reply.PieceHash,Reply.Name)
	//fmt.Println([]byte(Reply.PieceHash))
	if err!=nil{
		fmt.Println(err)
		return
	}
	file.Close()
}