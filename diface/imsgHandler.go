package diface

//消息管理抽象层

type IMsgHandle interface {
	//调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequset)
	//为消息添加具体的消息逻辑
	AddRouter(msgID uint32, router IRouter)

	// 启动一个Worker工作池 （开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池）
	StartWorkerPool()

	// 启动一个worker工作流程
	StartOneWorker(workerID int, taskQueue chan IRequset)

	// 将消息交给TaskQueue  由Worker进行处理
	SendMsgToTaskQueue(requset IRequset)
}
