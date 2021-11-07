package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//create a server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}
func (s *Server) Handler(conn net.Conn) {
	// current connect
	fmt.Println("connect successfully")
}

//start server
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
	}
	//close socket
	defer listener.Close()

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
