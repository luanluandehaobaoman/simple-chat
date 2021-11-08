package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn //对应user的socket/句柄
	server *Server
}

//创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

//用户上线的业务
func (u *User) Online() {
	//用户上线加到OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播用户上线
	u.server.BroadCast(u, "已上线")
}

//用户下线的业务
func (u *User) offline() {
	//用户下线从OnlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播用户下线
	u.server.BroadCast(u, "已下线")
}

//用户处理消息的业务
func (u *User) DoMessage(msg string) {
	u.server.BroadCast(u, msg)
}

//监听当前User channel的方法，一旦有消息就发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
