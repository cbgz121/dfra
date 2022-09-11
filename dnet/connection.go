package dnet

import (
	"dfra/diface"
	"dfra/pack"
	"dfra/untils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

// 连接模块
type Connection struct {
	TcpServer    diface.IServer
	conn         *net.TCPConn           //当前连接的 socket TCP套接字
	connID       uint32                 //连接的ID
	isClosed     bool                   //当前连接的状态
	ExitChan     chan bool              //告知当前连接已经退出的/停止的 channel (由Reader告知Writer退出)
	msgChan      chan []byte            //无缓冲的管道  用于读，写Groutine之间的消息通信
	MsgHandler   diface.IMsgHandle      //消息管理MsgID 和对应的业务处理api关系
	property     map[string]interface{} //连接属性集合
	propertyLock sync.RWMutex           //保护连接属性的锁
	active       uint32                 //说明连接是正常活跃的
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID =", c.connID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//创建一个拆包解包对象
		dp := pack.NewDataPack()

		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error 111", err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error  222", err)
				break
			}
		}
		msg.SetData(data)

		req := NewRequest(c, msg)

		// 表明是心跳包
		if msg.GetMsgId() == untils.GlobalObject.HeartBeatPackageId {
			c.SetActive(0)
		} else {
			if untils.GlobalObject.WorkerPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(req)
			} else {
				go c.MsgHandler.DoMsgHandler(req)
			}
		}
	}
}

// 写消息Groutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Gortinue is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler diface.IMsgHandle, server diface.IServer) diface.IConnection {
	return &Connection{
		conn:       conn,
		connID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
		TcpServer:  server,
		property:   make(map[string]interface{}),
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Sart()... ConnID=", c.conn)
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallOnConnStart(c)

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID=", c.connID)

	//如果当前连接已关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.conn.Close()

	//告知Witer关闭
	c.ExitChan <- true
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) SetActive(i uint32) {
	atomic.StoreUint32(&c.active, i)
}

func (c *Connection) GetActive() uint32 {
	return atomic.LoadUint32(&c.active)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) GetConnID() uint32 {
	return c.connID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) GetTcpServer() diface.IServer {
	return c.TcpServer
}

// 提供一个SendMsg方法  将我们要发送给客户端的数据，先进行封包 再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包， MsgDataLen|MsgId|Data
	dp := pack.NewDataPack()

	binaryMsg, err := dp.Pack(pack.NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
