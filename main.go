package main

import (
	"fmt"
	"go-ros-fog/conf"
	"go-ros-fog/server"
	"go-ros-fog/ziface"
	"go-ros-fog/znet"
)

// Test tiny ZINX
// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test tiny ZINX
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle..")
	// 先读取客户端的数据，再回写 Ping...Ping...Ping...
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("Ping...Ping...Ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 从配置文件读取配置
	conf.Init()

	// 创建server句柄，使用 tinyZinx api
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	go s.Serve()
	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}
