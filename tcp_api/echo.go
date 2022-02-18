package tcp_api

import (
	"fmt"
	"go-ros-fog/model"
	"go-ros-fog/ziface"
	"go-ros-fog/znet"
)

// EchoRouter 自定义路由
type EchoRouter struct {
	znet.BaseRouter
}

// echo 测试连通性
func (e *EchoRouter) Handle(request ziface.IRequest) {
	s := model.EchoMsg{}
	s.Unmarshal(request.GetData())
	err := request.GetConnection().SendMsg(0, s.Marshal())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
}
