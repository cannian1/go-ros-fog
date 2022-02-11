package main

import (
	"go-ros-fog/conf"
	"go-ros-fog/handle_tcp"
	"go-ros-fog/server"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	handle_tcp.TinyZinxServer()

	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}
