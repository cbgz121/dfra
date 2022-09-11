package dnet

import (
	"dfra/diface"
	"dfra/untils"
	"fmt"
	"log"
	"net"
	"time"
)

// IServer 接口的实现， 定义一个Server的服务器模块
type Server struct {
	Name           string                        //服务器的名称
	IPVersion      string                        //服务器绑定的ip版本
	IP             string                        //服务器监听的IP
	Port           int                           //服务器监听的端口
	MsgHandler     diface.IMsgHandle             //当前的Server添加的一个router，server注册的连接对应的处理业务
	ConnMgr        diface.IConnManager           //该server的连接管理器
	OnConnStart    func(conn diface.IConnection) //该Server创建连接后的Hook（钩子）方法
	OnConnStop     func(conn diface.IConnection) //该Server销毁之前的Hook（钩子）方法
	HeartBExitChan chan struct{}
}

// 初始化一个Server模块的方法
func NewServer(name string) diface.IServer {
	return &Server{
		Name:           untils.GlobalObject.Name,
		IPVersion:      "tcp4", //TODO  tcp6
		IP:             untils.GlobalObject.Host,
		Port:           untils.GlobalObject.TcpPort,
		MsgHandler:     NewMsgHandle(),
		ConnMgr:        NewConnManager(),
		HeartBExitChan: make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	fmt.Printf("from conf json: name:%s, IP:%s, Port:%v, \n", untils.GlobalObject.Name, untils.GlobalObject.Host, untils.GlobalObject.TcpPort)

	// 开启消息队列及worker工作池
	s.MsgHandler.StartWorkerPool()

	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error:", err)
		return
	}

	listenner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen", s.IPVersion, "err", err)
		return
	}

	fmt.Println("start dfra server success:", s.Name)
	var cID uint32
	cID = 0

	for {
		conn, err := listenner.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err", err)
			continue
		}

		//设置最大连接个数的判断  如果超出最大连接  那么则关闭此新的连接
		if s.ConnMgr.Len() >= untils.GlobalObject.MaxConn {
			conn.Close()
			continue
		}

		dealConn := NewConnection(conn, cID, s.MsgHandler, s)
		s.ConnMgr.Add(dealConn)
		cID++

		//启动当前的链接业务处理
		go dealConn.Start()
	}

}

func (s *Server) Stop() {
	fmt.Println("[STOP] dfra Server name", s.Name)
	//server 停止  清除所有的连接
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	go s.heartBeat()
	s.Start()
}

func (s *Server) AddRouter(msgID uint32, router diface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router succ~~")
}

func (s *Server) GetConnMgr() diface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection diface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection diface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn diface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn diface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}

// 服务端心跳检测，每5s将所有连接的active加1
func (s *Server) heartBeat() {
	log.Println("[server heart beat start SUCCESS]")

OverHeartBeat:
	for {
		time.Sleep(untils.GlobalObject.HeartRateInSecond)
		for id := range s.ConnMgr.GetConns() {
			conn, err := s.ConnMgr.Get(id)
			if err != nil {
				log.Println("[get nil connection]")
				continue
			}
			if conn.GetActive() == untils.GlobalObject.HeartFreshLevel {
				log.Println("[connection", conn.GetConnID(), "expired]")
				conn.Stop()
			} else {
				conn.SetActive(conn.GetActive() + 1)
			}
		}
		select {
		case <-s.HeartBExitChan:
			break OverHeartBeat
		default:
		}
	}
}
