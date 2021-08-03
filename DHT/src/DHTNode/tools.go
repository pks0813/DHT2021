package DHTNode

import (
	"crypto/sha1"
	"math/big"
	"net"
	"net/rpc"
	"time"
)
var DialTime=0
var FindSuccessorTime=0
type KVPair struct{
	Key string
	Value string
}
var Mod=new(big.Int).Exp(big.NewInt(2),big.NewInt(int64(M)),nil)
func PKSHash(tmp string) *big.Int {
	h:=sha1.New()
	h.Write([]byte(tmp))
	Q:=big.NewInt(0)
	Q.Mod(new(big.Int).SetBytes(h.Sum(nil)),Mod)
	return Q
}
func  PKSMerge(X *map[string]string,Y *map[string]string)  {
	for S1,S2:=range *Y {
		(*X)[S1]=S2
	}
}
func Inrange(L *big.Int,R *big.Int,X *big.Int)bool{
	//fmt.Println("L:",*L)
	//fmt.Println("R:",*R)
	//fmt.Println("X:",*X)
	FR:=big.NewInt(0)
	FR.Set(R)
	//fmt.Println("FR:",*FR)
	//fmt.Println(FR.Cmp(L))
	//fmt.Println(Mod)
	if FR.Cmp(L)<=0{
		FR.Add(R,Mod)
		//fmt.Println("NewFR:",FR)
	}
	//fmt.Println("L:",*L)
	//fmt.Println("FR:",*FR)
	//fmt.Println("X:",*X)
	FX:=big.NewInt(0)
	FX.Set(X)
	if FX.Cmp(L)<=0{
		FX.Add(X,Mod)
	}
	//fmt.Println(*L)
	//fmt.Println(*FR)
	//fmt.Println(*FX)
	return FX.Cmp(FR)<=0
}

func Dial(IP string) (*rpc.Client,error){
	//fmt.Println("Try to Dial",IP)
	var err error
	for i:=1;i<=5;i++ {
		NowCilent,err1:=rpc.Dial("tcp",IP)
		err=err1
		if err==nil {
			DialTime++
			return NowCilent,nil
		}
		time.Sleep(WaitTime)
	}
	return nil,err
}
func (this *Node)ping(IP string) bool{
	//fmt.Println("Try to Ping",IP)
	//var err1 error
	for i:=1;i<=5;i++ {
		//fmt.Println("Ping time:",i)
		Client,err:=rpc.Dial("tcp",IP)
		//err1=err
		if err==nil {
			_=Client.Close()
			return true
		}
		//time.Sleep(WaitTime)
	}
	//fmt.Println("Ping Fail",IP,err1)
	return false
}

func GetLocalAddress() string {
	var localaddress string

	ifaces, err := net.Interfaces()
	if err != nil {
		panic("init: failed to find network interfaces")
	}

	// find the first non-loopback interface with an IP address
	for _, elt := range ifaces {
		if elt.Flags&net.FlagLoopback == 0 && elt.Flags&net.FlagUp != 0 {
			addrs, err := elt.Addrs()
			if err != nil {
				panic("init: failed to get addresses for network interface")
			}

			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok {
					if ip4 := ipnet.IP.To4(); len(ip4) == net.IPv4len {
						localaddress = ip4.String()
						break
					}
				}
			}
		}
	}
	if localaddress == "" {
		panic("init: failed to find non-loopback interface with valid address on this node")
	}

	return localaddress
}
