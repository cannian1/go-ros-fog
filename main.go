package main

import (
	"go-ros-fog/conf"
	"go-ros-fog/ros"
	"go-ros-fog/tcp_main"
	"go-ros-fog/ticker"

	"go-ros-fog/server"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	tcp_main.TinyZinxServer()

	go ros.TopicServer()

	ticker.HandleAllTicker()

	// 装载路由
	r := server.NewRouter()

	r.Run(":3000")
}
