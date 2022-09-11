package dnet

import (
	"dfra/diface"
	"errors"
	"fmt"
	"sync"
)

//连接管理模块

type ConnManager struct {
	connections map[uint32]diface.IConnection //管理的连接集合
	connLock    sync.RWMutex
}

func NewConnManager() diface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]diface.IConnection),
	}
}

func (cm *ConnManager) GetConns() map[uint32]diface.IConnection {
	return cm.connections
}

// 添加连接
func (cm *ConnManager) Add(conn diface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn

	fmt.Println("connID=", conn.GetConnID(), "connection add to ConManager successfully :conn num =", cm.Len())
}

// 删除连接
func (cm *ConnManager) Remove(conn diface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除连接信息
	delete(cm.connections, conn.GetConnID())

	fmt.Println("connID=", conn.GetConnID(), "remove to ConManager successfully :conn num =", cm.Len())
}

// 根据connID获取连接
func (cm *ConnManager) Get(connID uint32) (diface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		//找到了
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}
}

// 得到当前连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清除并终止所有连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All connections succ ! conn num =", cm.Len())
}

// ClearOneConn  利用ConnID获取一个链接 并且删除
func (cm *ConnManager) ClearOneConn(connID uint32) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	connections := cm.connections
	if conn, ok := connections[connID]; ok {
		conn.Stop()
		delete(connections, connID)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
		return
	}

	fmt.Println("Clear Connections ID:  ", connID, "err")
	return
}
