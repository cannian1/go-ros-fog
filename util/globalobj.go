package util

import (
	"go-ros-fog/ziface"
	"encoding/json"
	"os"
)

// 存储一切有关 tinyZinx 框架的全局参数，供其他模块使用
// 一些参数可以通过 tinyZinx.json 由用户进行配置

type GlobalObj struct {
	// Server
	TcpServer ziface.IServer // 当前 tinyZinx 全局的Server对象
	Name      string         // 当前服务器的名称
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	// tinyZinx
	Version           string // 当前 tinyZinx 的版本号
	MaxConn           int    // 当前服务器主机允许的最大连接数
	MaxPackageSize    uint32 // 当前 tinyZinx 框架数据包的最大值
	WorkerPoolSize    uint32 // 当前业务工作 Worker池 的最大 goroutine 数量
	MaxWorkerTaskLen  uint32 // tinyZinx 框架允许用户最多开辟多少个 Worker （限定条件）
	HeartbeatInterval int    // 连接的心跳间隔
}

// GlobalObject 定义一个全局对外 GlobalObj
var GlobalObject *GlobalObj

// ReloadFromConf 从 tinyZinx.json 加载用于用户自定义的参数
func (g *GlobalObj) ReloadFromConf() {
	// go 1.16 弃用此包 ioutil.ReadFile()
	data, err := os.ReadFile("conf/tinyzinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// ReloadFromEnv 从 .env 里加载参数
func (g *GlobalObj) ReloadFromEnv(){
	// TODO 还没写
}

// 初始化当前的 GlobalObject
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:              "tinyZinxServerApp",
		Host:              "0.0.0.0",
		TcpPort:           8999,
		Version:           "V0.11", // 这里是写死的，不该被客户端配置文件更改
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,   // worker 工作池的队列个数
		MaxWorkerTaskLen:  1024, // 每个 worker 对应的消息队列任务的数量最大值
		HeartbeatInterval: 60,
	}

	// 尝试从 conf/tinyzinx.json去加载一些用户自定义的参数
	GlobalObject.ReloadFromConf()
}
