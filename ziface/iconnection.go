package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	// Start 启动连接 让当前的连接准备开始工作
	Start()
	// Stop 停止连接 结束当前连接的工作
	Stop()
	// GetTCPConnection 获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的 TCP状态 IP PORT
	RemoteAddr() net.Addr
	// SendMsg 发送数据，将数据发送给远程客户端
	SendMsg(msgId uint32, data []byte) error

	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)

	// HeartBeatChecker 检查心跳是否健康
	HeartBeatChecker()
	// IsAlive 心跳状态
	IsAlive() bool
	// KeepAlive 更新心跳
	KeepAlive()
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
