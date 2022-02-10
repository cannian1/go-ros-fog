package znet

import (
	"fmt"
	"go-ros-fog/util"
	"go-ros-fog/ziface"
	"net"
)

// Server IServer的接口实现
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// v0.6 当前server 的消息管理模块，用来绑定 MsgID 和 对应的处理业务 API 关系
	MsgHandler ziface.IMsgHandle
	// 该 server的连接管理器
	ConnMgr ziface.IConnManager
	// server 创建连接之后自动调用的 Hook 函数 (广播之类的业务)
	OnConnStart func(conn ziface.IConnection)
	// server 销毁连接之后自动调用的 Hook 函数
	OnConnStop func(conn ziface.IConnection)
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[tinyZinx] Server Name: %s, listening at IP: %d is starting\n", util.GlobalObject.Host, util.GlobalObject.TcpPort)
	fmt.Printf("[tinyZinx] Version %s , MaxConn:%d ,MaxPackageSize:%d\n", util.GlobalObject.Version, util.GlobalObject.MaxConn, util.GlobalObject.MaxPackageSize)

	// 异步
	go func() {
		// 开启消息队列及 worker 工作池
		s.MsgHandler.StartWorkerPool()

		// 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error :", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start tinyZinx server succ,", s.Name, "succ, Listening")

		// v0.2 分配连接id
		var cid uint32
		cid = 0

		// 阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 设置最大连接个数的判断，如果超过最大连接，则关闭此连接
			if s.ConnMgr.Count() >= util.GlobalObject.MaxConn {
				// TODO 给客户端返回一个超出最大连接的错误包
				// fmt.Println("=======>Too Many Connections, MaxConn =  ", util.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// v0.2 调用Connection模块处理连接业务
			// 将处理新连接的业务方法和 conn 进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// v0.2 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

// Serve 运行服务器
func (s *Server) Serve() {
	s.Start()
	// TODO 可以做一些启动服务器之后的额外工作
	// 阻塞状态
	select {}
}

// GetConnMgr 获取当前server的连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succeed!!")
}

// NewServer 初始化Server模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       util.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         util.GlobalObject.Host,
		Port:       util.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// SetOnConnStart 注册 OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册 OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用 OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用 OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}

}
