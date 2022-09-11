package main

import (
	"dfra/diface"
	"dfra/dnet"
	"fmt"
)

// 基于dfra框架开发的服务器端应用程序

// ping test 自定义路由
type PingRouter struct {
	dnet.BaseRouter
}

// Test Handle
func (b *PingRouter) Handle(request diface.IRequset) {
	fmt.Println("Call ping Router Handle...")
	//先读取客户端的数据  再写回ping...ping...ping
	fmt.Println("recv from client:msgID =", request.GetMsgId(), ",data =", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// hello dfra test 自定义路由
type HellodfraRouter struct {
	dnet.BaseRouter
}

func (b *HellodfraRouter) Handle(request diface.IRequset) {
	fmt.Println("Call HellodfraRouter Router Handle...")
	//先读取客户端的数据  再写回
	fmt.Println("recv from client:msgID =", request.GetMsgId(), ",data =", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello welcome dfra"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建一个server句柄，使用dfra的api
	s := dnet.NewServer("[dfra]")

	//给当前dfra框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HellodfraRouter{})

	//启动server
	s.Serve()
}
