package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User //在线用户的列表
	mapLock   sync.RWMutex
	Message   chan string //消息广播的channel
}

//create a server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听广播消息的channel ：Message，一旦有消息就发送给全部在线的User
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		//将msg发给所有在线的User
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

//广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// current connect
	fmt.Println("connect successfully from:  ", conn.RemoteAddr().String())
	user := NewUser(conn, s)
	//用户上线加到OnlineMap

	user.Online()

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//用户下线
				user.offline()
			}
			if err != nil && err != nil {
				fmt.Println("conn Read err", err)
			}
			//提取用户消息(去除'\n')
			msg := string(buf[:n-1])

			//广播消息
			user.DoMessage(msg)
		}
	}()

	//当前handler阻塞
	//select {}
}

//启动服务器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
	}
	//close socket
	defer listener.Close()

	//启动监听Message的goroutine
	go s.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		//do handler
		go s.Handler(conn)
	}
}
