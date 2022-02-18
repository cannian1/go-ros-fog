package znet

import (
	"go-ros-fog/util"
	"go-ros-fog/ziface"
	"fmt"
	"strconv"
)

// 消息处理模块实现

type MsgHandle struct {
	// 存放每个 MsgID 对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责 Worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作 Worker池的 Worker 数量
	WorkerPoolSize uint32
}

// NewMsgHandle 初始化创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: util.GlobalObject.WorkerPoolSize, // 从全局配置参数中获取
		TaskQueue:      make([]chan ziface.IRequest, util.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度/执行对应 Router 消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 从 request 中找到 msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! NEED register!")
		return
	}
	// 根据 MsgID 调度对应 router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断 当前 msg 绑定的Api处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id 已经注册了
		panic("repeat api , msgID =" + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success!")
}

// StartWorkerPool 启动一个 Worker 工作池
// (开启工作池只能发生一次，只有一个工作池)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据 workerPoolSize 分别开启 Worker，每个Worker用一个 go 取
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个 worker被启动
		// 当前 worker 对应的 channel 消息队列，开辟空间 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, util.GlobalObject.MaxWorkerTaskLen)

		// 启动当前的 Worker，阻塞等待消息从 channel 传入
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个 Worker 工作流程
func (mh *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started")

	// 不断地阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息传入，出列的是一个客户端的 Request，执行当前 Request 绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToQueue 将消息交给TaskQueue，由 Worker处理
func (mh *MsgHandle) SendMsgToQueue(request ziface.IRequest) {
	// 将消息平均分配给不同的 worker
	// 如果是分布式的负载均衡可以考虑别的负载均衡算法以及根据对方 IP 地域分配
	// 根据客户端建立的 ConnID 进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgID(), " to WorkerID = ", workerID)

	// 将消息发送给队员的 Worker 的TaskQueue
	mh.TaskQueue[workerID] <- request
}
