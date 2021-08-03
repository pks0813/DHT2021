####DHT文件
#####package DHTNode
DHTNode.go:实现了dhtnode接口所需要的函数实现  
Node.go：实现了DHT内部节点所需要的函数并且实现节点之间的信号拨打  
tools.go:实现了DHTNode.go与Node.go所需要的额外函数
#####package main
userdef.go:实现了dhtNode NewNode(int)函数
###application：File Manage System
在DHT的基础上实现了package FMS 并且使用了Git上bencode库 对种子文件进行读写    
在同一局域网不同的机器或同一个机器不同的端口可进行文件传输  
且利用go所支持的并行对文件进行拆分 加快文件下载速度  
初始会在当前目录下创建Download Upload Seed 三个文件夹  
将传输文件放在Upload里 输出得到的种子会被写入Seed文件夹的指定文件里  
下载时会读取Seed文件夹中指定的文件 输出至Download文件夹里
