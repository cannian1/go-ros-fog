package tcp_api

import (
	"fmt"
	"go-ros-fog/tcp_core"
	"go-ros-fog/tcp_model"
	"go-ros-fog/ziface"
	"go-ros-fog/znet"
)

type DealWithIoTData struct {
	znet.BaseRouter
}

// Handle 处理连接业务
func (dwid *DealWithIoTData) Handle(request ziface.IRequest) {
	fmt.Println("***处理连接业务***")
	edgeDevice := tcp_core.EDMgrObj.GetEdgeDeviceByDId(request.GetConnection().GetConnID())
	data := request.GetData()
	// json.Unmarshal(data, edgeDevice.Sensor)
	t := edgeDevice.Marshal2Sensor(data)
	// 别句柄了，句柄会共享
	resMsg := &tcp_model.ResRelay{}

	// TODO: 数据持久化
	// TODO: 越界报警与自动处置
	// TODO: 可视化

	switch {
	case t.Temperature > 50:
		fmt.Println(t.Temperature, "温度超标了")
		resMsg.Temperature = true
		fallthrough
	case t.LightLevel > 50:
		fmt.Println(t.LightLevel, "光强超标了")
		resMsg.LightLevel = true
		fallthrough
	case t.Smog > 50:
		fmt.Println(t.Smog, "粉尘浓度超标")
		resMsg.Smog = true
		err := request.GetConnection().SendMsg(15001, resMsg.Marshal())
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("一切正常", t)
		err := request.GetConnection().SendMsg(15000, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}
