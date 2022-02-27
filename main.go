package main

import (
	"go-ros-fog/conf"
	"go-ros-fog/ros"
	"go-ros-fog/tcp_main"

	"go-ros-fog/server"
)

func main() {
	// 从配置文件读取配置
	conf.Init()
	// t:=ros.GetNodes()
	// fmt.Println("***ROS NODES*** ",t)

	tcp_main.TinyZinxServer()
	ros.TopicServer()

	// t := ros.GetNodes()
	// fmt.Println(t)
	// t1 := ros.BusinessNode{}
	// time.Sleep(1 * time.Second)
	// go t1.InitSubscriber()
	// go t1.InitPublisher()

	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}
