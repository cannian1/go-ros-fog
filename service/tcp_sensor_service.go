package service

import (
	"encoding/json"
	"fmt"
	"go-ros-fog/cache"
	"go-ros-fog/model"
	"go-ros-fog/serializer"
	"strconv"
)

type TcpSensorService struct {
}

func (service *TcpSensorService) SensorValue(id string) serializer.Response {
	// id is unused 先接着，设备多了可以塞进去拼接redis
	err := cache.RedisClient.HGetAll(cache.SensorValue).Err()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Msg:   "redis 数据库连接出错",
			Error: err.Error(),
		}
	}

	data := serializer.BuildSensor()
	if len(data) == 0 {
		return serializer.Response{
			Data: data,
			Msg:  "The sensor is not ready or disconnected",
		}
	}
	return serializer.Response{
		Data: data,
	}
}

// TODO: 每个设备建立连接时向redis注册deviceid，从里面读出所有的id，遍历hash
// func (service *TcpSensorService)SensorValues()serializer.Response{

// }

func (service *TcpSensorService) SetSensorThreshold(id string, body []byte) serializer.Response {
	// id is unused 先接着，设备多了可以塞进去拼接redis

	id = "1"
	db_primary, _ := strconv.Atoi(id)
	fmt.Println(db_primary)
	// 反序列化body到map
	recv := make(map[string]interface{})
	if err := json.Unmarshal(body, &recv); err != nil {
		return serializer.Response{
			Code:  serializer.CodeParamErr,
			Msg:   "body格式错误",
			Error: err.Error(),
		}
	}

	// 过滤有效字段
	redisSave := make(map[string]interface{})
	redisSave["equipment_id"] = id
	redisSave["temperature"] = recv["temperature"]
	redisSave["light_intensity"] = recv["light_intensity"]
	redisSave["smog"] = recv["smog"]

	// 存入redis
	redisErr := cache.RedisClient.HMSet(cache.SensorThreshold+id, redisSave).Err()
	if redisErr != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Msg:   "redis 连接异常",
			Error: redisErr.Error(),
		}
	}

	var sensorThreshold model.SensorsThreshold
	if err := json.Unmarshal(body, &sensorThreshold); err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Msg:   "反序列化到model失败，一般这种问题是空指针",
			Error: err.Error(),
		}
	}

	sensorThreshold.EquipmentID = uint32(db_primary)
	// Save 会保存所有的字段，即使字段是零值
	err := model.DB.Save(&sensorThreshold).Error
	if err != nil {
		// 数据库还没建表
		createErr := model.DB.Create(&sensorThreshold).Error
		if createErr != nil {
			return serializer.Response{
				Code:  50001,
				Msg:   "保存阈值失败",
				Error: createErr.Error(),
			}
		}
	}

	return serializer.Response{
		Data: true,
	}
}

// TODO: 登陆以后先查redis，为空就去mysql里查，写入redis