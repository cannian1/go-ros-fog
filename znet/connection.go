package znet

import (
	"errors"
	"fmt"
	"go-ros-fog/util"
	"go-ros-fog/ziface"
	"io"
	"net"
	"sync"
	"time"
)

// Connection 连接模块
type Connection struct {
	// 当前 Conn 隶属于哪个 server
	TcpServer ziface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool

	// 告知当前连接已经退出/停止 channel,(由 Reader告知 Writer 退出)
	ExitChan chan bool

	// 无缓冲管道，用于读、写 goroutine
	msgChan chan []byte

	// 消息的管理 MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex

	// 查看、更新心跳状态要加锁
	sync.RWMutex
	// 最近一次心跳时间
	lastHeartbeatTime time.Time
}

// NewConnection 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:         server,
		Conn:              conn,
		ConnID:            connID,
		MsgHandler:        handle,
		isClosed:          false,
		msgChan:           make(chan []byte),
		ExitChan:          make(chan bool, 1),
		property:          make(map[string]interface{}),
		lastHeartbeatTime: time.Now(),
	}
	// 将 conn 加入到 ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running..]")
	defer fmt.Println("[conn Reader is exited!] connID = ", c.ConnID, " remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//// 读客户端的数据到buf中
		//buf := make([]byte, util.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("receive buf err", err)
		//	continue
		//}

		// 保持心跳
		c.KeepAlive()

		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的 Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		// 拆包，得到 msgID 和 msgDataLen 放在 msg 消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}

		// 根据 dataLen,再次读取 Data,放到 msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 判断是否开启工作池
		if util.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给 Worker 工作池处理即可
			c.MsgHandler.SendMsgToQueue(&req)
		} else {
			// v0.3 不需要 handleApi 了，v0.6 删掉单路由
			// 从路由中找到注册绑定的Conn对应的router调用
			// 根据绑定好的 MsgID 找到对应处理 api业务执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// StartWriter 写消息的 goroutine，专门发送给客户端消息的模块
// 读写分离可以在发数据前进行处理
func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Println("[conn Writer exit!] ", c.RemoteAddr().String())

	// 不断地阻塞等待 channel 的消息，写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			// 代表 Reader 已经退出，此时 Writer 也要退出
			return
		}
	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {

	fmt.Println("Conn Start()... ConnId = ", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
	go c.StartWriter()

	// 按照开发者传递进来的，创建连接之后需要调用的处理业务，执行对应 Hook 函数
	c.TcpServer.CallOnConnStart(c)
	c.HeartBeatChecker()
}

// Stop 停止连接 结束当前连接的工作
// 任何一个出错都会 break，执行 defer 中的 Stop
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed{
		return
	}
	c.isClosed = true

	// 调用开发者注册的 销毁连接之前需要执行的业务 Hook函数
	c.TcpServer.CallOnConnStop(c)
	// 关闭 socket 连接
	_ = c.Conn.Close()

	// 告知 Writer 关闭
	c.ExitChan <- true

	// 将当前连接从 ConnMgr 中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的 TCP状态 IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg  提供一个 SendMsg 方法 将我们要发给客户端的数据先进行封包，再发送
// 在用户使用时的自定义回调函数（main server.go）使用
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("the connection is closed when the message is sent")
	}

	// 将 data 进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("pack err msg")
	}

	// 将数据发送给客户端,V0.7 发送到 StartWriter 中的管道处理
	// 把消息打包完通过 channel 送给 writer
	c.msgChan <- binaryMsg

	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加连接属性
	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// HeartBeatChecker 检查心跳是否健康
func (c *Connection) HeartBeatChecker() {
	var timer *time.Timer
	timer = time.NewTimer(time.Duration(util.GlobalObject.HeartbeatInterval) * time.Second)
	for {
		select {
		case <-timer.C:
			if !c.IsAlive() {
				c.Stop()
				goto EXIT
			}
			timer.Reset(time.Duration(util.GlobalObject.HeartbeatInterval) * time.Second)
		case <-c.ExitChan:
			timer.Stop()
			goto EXIT
		}
	}
EXIT:
	// 确保连接被关闭
}

// IsAlive 心跳状态
func (c *Connection) IsAlive() bool {
	now := time.Now()
	c.Lock()
	defer c.Unlock()
	// 连接被关闭 或者 现在时刻离上次心跳时间超过设定值
	if c.isClosed || now.Sub(c.lastHeartbeatTime) > time.Duration(util.GlobalObject.HeartbeatInterval)*time.Second {
		return false
	}
	return true
}

// KeepAlive 更新心跳
func (c *Connection) KeepAlive() {
	now := time.Now()
	c.Lock()
	defer c.Unlock()

	c.lastHeartbeatTime = now
}
