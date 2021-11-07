package main

type Server struct {
	Ip   string
	Prot int
}

//create a server
func NewServer(ip string, port int) *Server {
	server:=&Server{
		Ip: ip,
		Prot: port,
	}
	return server
}

//start server
func (S *Server)Start()  {
	//socket listen
	listener,err :=  net.Listen()

	//accept

	//do handler

	//close socket 
}


