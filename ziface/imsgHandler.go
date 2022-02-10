package ziface

// 消息管理抽象层

type IMsgHandle interface {
	// DoMsgHandler 调度/执行对应 Router 消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动一个 Worker 工作池
	StartWorkerPool()
	// SendMsgToQueue 将消息交给TaskQueue，由 Worker处理
	SendMsgToQueue(request IRequest)
}
