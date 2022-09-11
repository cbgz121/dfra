package dnet

import (
	"dfra/diface"
	"dfra/untils"
	"fmt"
	"sync"
)

//消息处理模块的实现

type MsgHandle struct {
	Apis           sync.Map               //存放每个MsgID所对应的处理方法
	TaskQueue      []chan diface.IRequset //负责Worker取任务的消息队列
	WorkerPoolSize uint32                 //业务工作Worker池的worker数量
}

// 初始化/创建MsgHandle 方法
func NewMsgHandle() diface.IMsgHandle {
	return &MsgHandle{
		WorkerPoolSize: untils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan diface.IRequset, untils.GlobalObject.WorkerPoolSize),
	}
}

// 调度/执行对应的Router消息处理方法
func (m *MsgHandle) DoMsgHandler(request diface.IRequset) {
	handler1, ok := m.Apis.Load(request.GetMsgId())
	if !ok {
		fmt.Println("api msgID =", request.GetMsgId(), "is NOT FOUND! Need Register!")
		return
	}
	//根据MsgID 调度对应的router业务即可
	handler := handler1.(diface.IRouter)
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的消息逻辑
func (m *MsgHandle) AddRouter(msgID uint32, router diface.IRouter) {
	if _, ok := m.Apis.Load(msgID); ok {
		fmt.Println("repeat api, msgID", msgID)
		return
	}
	//添加msg与api的绑定关系
	m.Apis.Store(msgID, router)
	fmt.Println("Add api MsgID =", msgID, "succ!")
}

// 启动一个Worker工作池 （开启工作池的动作只能发生一次，一个dfra框架只能有一个worker工作池）
func (m *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize分别开启worker  每个worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan diface.IRequset, untils.GlobalObject.WorkerPoolSize)
		//启动当前的worker  阻塞等待消息从channel 传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan diface.IRequset) {
	fmt.Println("Worker ID =", workerID, "is started...")

	//不断阻塞等待对应的消息队列的消息
	for {
		select {
		//如果有消息过来， 出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue  由Worker进行处理
func (m *MsgHandle) SendMsgToTaskQueue(requset diface.IRequset) {
	//根据消息平均分配给不同的worker
	//根据客户端建立的connID来进行分配
	// TODO 消息均衡算法的优化
	workerID := requset.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID=", requset.GetConnection().GetConnID(), "request MsgID=", requset.GetMsgId(), "to workerID=", workerID)

	//将消息发送给对应的worker的TaskQueue即可
	m.TaskQueue[workerID] <- requset
}
