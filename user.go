package main

import (
	"fmt"
	"net"
	"strings"
)

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

//给当前User对应的客户端发消息
func (u *User) SendMsg(msg string) {
	_, err := u.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
	}
}

//用户处理消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户有哪些
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线..\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()

		//支持改名
	} else if len(msg) > 7 && msg[:7] == "rename|" { //消息格式：rename|Arthur
		newName := strings.Split(msg, "|")[1]

		//判断name是否已经存在
		_, ok := u.server.OnlineMap[newName]

		if ok {
			u.SendMsg("当前用户名已存在")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("您已经更新用户名为：" + u.Name + "\n")
		}

		//支持私聊，消息格式：to|张三|消息内容
	} else if len(msg) > 4 && msg[:3] == "to|" {

		//1.获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("消息格式不正确，请使用'to|张三|消息内容'\n")
			return
		}

		//2.根据用户名获取User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("用户" + remoteName + "不存在\n")
			return
		}

		//3.获取消息内容并发送
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("无消息内容，请重新发送\n")
			return
		}
		remoteUser.SendMsg("来自" + u.Name + "的消息:" + content + "\n")

	} else {

		u.server.BroadCast(u, msg)
	}
}

//监听当前User channel的方法，一旦有消息就发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}
