package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	conn       net.Conn
	name       string
}

func NewClient(serverip string, serverport int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverip,
		ServerPort: serverport,
	}
	//连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverip, serverport))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn
	//返回对象
	return client
}

var serverIp string
var serverPort int

//./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>连接服务器失败……")
	}
	fmt.Println(">>>>>>>连接服务器成功……")
	select {}
}
