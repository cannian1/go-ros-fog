package tcp_api

import (
	"fmt"
	"go-ros-fog/cache"
	"go-ros-fog/tcp_core"
	"go-ros-fog/tcp_model"
	"go-ros-fog/ziface"
	"go-ros-fog/znet"
	"reflect"
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

	// 反射指向结构体指针 (.Elem)
	ptrToTypeElemSe := reflect.TypeOf(t).Elem()
	ptrToValueElemSe := reflect.ValueOf(t).Elem()
	// 获取结构体有几个字段
	seNumField := ptrToTypeElemSe.NumField()

	// 把结构体转为 map
	sensorMap := make(map[string]interface{})
	for i := 0; i < seNumField; i++ {
		// map 的 key 设为结构体成员的字段名，value设为成员的值
		// 最后一定要.interface()，否则map的value全是reflect.value类型，redis无法序列化
		sensorMap[ptrToTypeElemSe.Field(i).Name] = ptrToValueElemSe.Field(i).Interface()
	}

	// for k, v := range sensorMap {
	// 	fmt.Println(k, reflect.TypeOf(v))
	// }

	// 存到 redis 数据库，方便跨语言共享
	err := cache.RedisClient.HMSet(cache.SensorValue, sensorMap).Err()
	if err != nil {
		panic(err)
		//fmt.Println(err)
	}
	// 别句柄了，句柄会共享
	resMsg := &tcp_model.ResRelay{}

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