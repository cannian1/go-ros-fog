package tcp_main

import (
	"fmt"
	"go-ros-fog/cache"
	"go-ros-fog/tcp_api"
	"go-ros-fog/tcp_core"
	"go-ros-fog/tcp_model"
	"go-ros-fog/ziface"
	"go-ros-fog/znet"
)

// 自定义 msgID
const (
	echo               = iota // 测试连通性
	DEAL_WITH_IOT_DATA = 3000 //
)

// DoConnectionBegin 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("执行连接开始工作...")
	edgeDevice := tcp_model.NewEdgeDevice(conn)
	// 同步 id 给客户端，服务端给其分配的 id
	edgeDevice.SyncDid()
	// 向设备管理器注册设备对象
	tcp_core.EDMgrObj.AddEdgeDevice(edgeDevice)

	// 在redis中自增连接数
	// tcp_conn:count
	err := cache.RedisClient.Incr(cache.TCPConnCount).Err()
	if err != nil {
		panic(err)
	}

	// // 设置一些连接属性
	// fmt.Println("Set Conn property")
	// conn.SetProperty("ServiceType", "TCP_SCM")
	// conn.SetProperty("Description", "earth")

}

// DoConnectionLost 连接断开之前需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("执行连接断开前工作...")
	id := conn.GetConnID()
	fmt.Println(id, "下线")

	// 从设备管理器删除设备对象
	tcp_core.EDMgrObj.Remove(id)

	// 释放 redis 保存该连接的内容,释放后再取只能取到空切片
	result, err := cache.RedisClient.HKeys(cache.SensorValue).Result()
	if err != nil {
		fmt.Println("Get conn_lost hkeys fail ", err)
	}

	err = cache.RedisClient.HDel(cache.SensorValue, result...).Err()
	if err != nil {
		fmt.Println("remove SensorValue err", err)
	}

	// 在redis中自减连接数
	err = cache.RedisClient.Decr(cache.TCPConnCount).Err()
	if err != nil {
		fmt.Println("remove TCPConnCount err", err)
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
	s.AddRouter(echo, &tcp_api.EchoRouter{})
	s.AddRouter(DEAL_WITH_IOT_DATA, &tcp_api.DealWithIoTData{})

	// 启动server
	go s.Serve()
}
