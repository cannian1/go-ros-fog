package handle_tcp

import (
	"fmt"
	"go-ros-fog/cache"
	"go-ros-fog/model"
	"go-ros-fog/ziface"
	"go-ros-fog/znet"
)

// 自定义 msgID
const (
	echo          = iota     // 测试连通性
	recvSensors   = iota + 1 // 接收传感器数据
	controlRelays            // 控制继电器
	controlServos            // 控制伺服电机
)

// EchoRouter 自定义路由
type EchoRouter struct {
	znet.BaseRouter
}

// SensorsRouter 自定义路由
type SensorsRouter struct {
	znet.BaseRouter
}

// RelayRouter 自定义路由
type RelayRouter struct {
	znet.BaseRouter
}

// ServoRouter 自定义路由
type ServoRouter struct {
	znet.BaseRouter
}

// echo 测试连通性
func (e *EchoRouter) Handle(request ziface.IRequest) {
	s := model.EchoMsg{}
	s.Unmarshal(request.GetData())
	err := request.GetConnection().SendMsg(echo+1, s.Marshal())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
}

// 处理传感器数据
func (sr *SensorsRouter) Handle(request ziface.IRequest) {
	s := model.T_Sensors{}
	s.Unmarshal(request.GetData())
	// TODO 数据持久化
	// TODO 越界报警与自动处置
	// TODO 可视化
	fmt.Println(s)
}

// 计算继电器工作时间的数据,以及获取继电器当前工作状态
func (rr *RelayRouter) Handle(request ziface.IRequest) {
	s := model.T_Relays{}
	s.Unmarshal(request.GetData())
	// TODO
	// 保持第一次 状态为 true的时间 和 第一次状态为 false 的时间，计算差值。既继电器工作的时间
	fmt.Println(s)
}

// 获取伺服电机状态和角度
func (sr *ServoRouter) Handle(request ziface.IRequest) {
	s := model.T_Servos{}
	s.Unmarshal(request.GetData())
	// TODO
	// 数据持久化
}

// DoConnectionBegin 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	// 将当前连接的ID输出
	fmt.Println(conn.GetConnID())
	
	// 在redis中自增连接数
	// tcp_conn:count
	err := cache.RedisClient.Incr(cache.TCPConnCount).Err()
	if err != nil {
		panic(err)
	}

	// 设置一些连接属性
	fmt.Println("Set Conn property")
	conn.SetProperty("ServiceType", "TCP_SCM")
	conn.SetProperty("Description", "earth")

	// 获取连接属性
	if serviceType, err := conn.GetProperty("ServiceType"); err == nil {
		// 在redis中 自增连接类型与连接数
		// tcp_conn:service_type:TCP_SCM
		err = cache.RedisClient.Incr(fmt.Sprintf("%s:%s", cache.TCPConnServiceType, serviceType)).Err()
		if err != nil {
			panic(err)
		}
	}
}

// DoConnectionLost 连接断开之前需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost..")

	// 在redis中自减连接数
	err := cache.RedisClient.Decr(cache.TCPConnCount).Err()
	if err != nil {
		fmt.Println("remove TCPConnCount err", err)
	}

	// 获取连接属性
	if serviceType, err := conn.GetProperty("ServiceType"); err == nil {
		temp := fmt.Sprintf("%s:%s", cache.TCPConnServiceType, serviceType)

		// 在redis中 自减连接类型与连接数
		// tcp_conn:service_type:TCP_SCM
		if err = cache.RedisClient.Decr(temp).Err(); err != nil {
			fmt.Println("remove ServiceType err", err)
		}
	}

}

// TinyZinxServer 提供TCP长连接服务
func TinyZinxServer() {
	// 创建server句柄，使用 tinyZinx api
	s := znet.NewServer()

	// 注册连接 Hook 钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 注册路由
	s.AddRouter(echo, &EchoRouter{})
	s.AddRouter(recvSensors, &SensorsRouter{})
	s.AddRouter(controlRelays, &RelayRouter{})
	s.AddRouter(controlServos, &ServoRouter{})

	// 启动server
	go s.Serve()
}
