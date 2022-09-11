package diface

import "net"

type IConnection interface {
	//启动连接  让当前的连接准备开始工作
	Start()

	//停止连接  结束当前的连接工作
	Stop()

	//获取当前的连接的绑定socket conn
	GetTCPConnection() *net.TCPConn

	//获取当前连接模块的连接ID
	GetConnID() uint32

	//获取远程客户端的TCP状态 IP, port
	RemoteAddr() net.Addr

	//发送数据  将数据返送给远程的客户端
	SendMsg(msgId uint32, data []byte) error

	//获得TcpServer
	//TODO 这是自己添加的
	GetTcpServer() IServer

	//设置连接属性
	SetProperty(key string, value interface{})

	//获取连接属性
	GetProperty(key string) (interface{}, error)

	//移除连接属性
	RemoveProperty(key string)

	SetActive(i uint32)

	GetActive() uint32
}
