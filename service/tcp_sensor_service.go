package service

import (
	"go-ros-fog/cache"
	"go-ros-fog/serializer"
)

type TcpSensorService struct {
}

func (service *TcpSensorService) SensorValue(id string) serializer.Response {
	// id is unused 先接着，设备多了可以塞进去拼接redis
	err := cache.RedisClient.HGetAll(cache.SensorValue).Err()
	if err != nil {
		return serializer.Response{
			Code:  50001,
			Msg:   "redis 数据库连接出错",
			Error: err.Error(),
		}
	}

	data := serializer.BuildSensor()
	if len(data) == 0 {
		return serializer.Response{
			Data: data,
			Msg: "The sensor is not ready or disconnected",
		}
	}
	return serializer.Response{
		Data: data,
	}
}
