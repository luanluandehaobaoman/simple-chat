package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	conn       net.Conn
	Name       string
	flag       int //当前客户端模式
}

func NewClient(serverip string, serverport int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverip,
		ServerPort: serverport,
		flag:       999,
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
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}

}

//处理server回应消息，知道显示到标准输出
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)

	/*	//等价于
		for {
			buf:=make([]string,4096)
			client.conn.Read(buf)
			fmt.Println(buf)
		}*/
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>>>>请输入聊天内容，exit退出")

	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>>>>请输入聊天内容，exit退出")

		fmt.Scanln(&chatMsg)
	}

}

//查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	client.SelectUsers()
	var remoteName string
	var chatMsg string

	fmt.Println(">>>>>>>请输入聊天对象[用户名]，exit退出：")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>>>>>请输入消息内容，exit退出：")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			for chatMsg != "exit" {
				if len(chatMsg) != 0 {
					sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
					_, err := client.conn.Write([]byte(sendMsg))
					if err != nil {
						fmt.Println("conn.Write err:", err)
						break
					}
				}
				chatMsg = ""
				fmt.Println(">>>>>>>>请输入聊天内容，exit退出")
				fmt.Scanln(&chatMsg)
			}
		}
		client.SelectUsers()
		fmt.Println(">>>>>>>请输入聊天对象[用户名]，exit退出：")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) UpdateNmae() {
	fmt.Println(">>>>>>>>请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write.err", err)
		//return false
	}
	//return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}

		//根据不通模式处理业务
		switch c.flag {
		case 1:
			//fmt.Println(">>>>>>>公聊模式选择")
			c.PublicChat()
			break
		case 2:
			//fmt.Println(">>>>>>>私聊模式选择")
			c.PrivateChat()
			break
		case 3:
			//fmt.Println(">>>>>>>更新用户名选选择")
			c.UpdateNmae()

			break
		}
	}
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>连接服务器失败……")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>>>>>连接服务器成功……")
	client.Run()
}
