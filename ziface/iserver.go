package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()

	// AddRouter v0.3 路由功能：给当前的服务注册一个路由方法，供客户端连接使用
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取当前server的连接管理器
	GetConnMgr() IConnManager

	// SetOnConnStart 注册 OnConnStart 钩子函数的方法
	SetOnConnStart(func(conn IConnection))
	// SetOnConnStop 注册 OnConnStop 钩子函数的方法
	SetOnConnStop(func(conn IConnection))
	// CallOnConnStart 调用 OnConnStart 钩子函数的方法
	CallOnConnStart(conn IConnection)
	// CallOnConnStop 调用 OnConnStop 钩子函数的方法
	CallOnConnStop(conn IConnection)
}
