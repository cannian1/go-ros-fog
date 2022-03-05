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
		// 自动迁移是在启动程序的时候进行的，运行后再删表就会在这里报错
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

//  登陆以后先查redis，为空就去mysql里查，由前端处理写入redis
func (service *TcpSensorService) GetSensorThreshold(id string) serializer.Response {

	id = "1"
	result, err := cache.RedisClient.HGetAll(cache.SensorThreshold + id).Result()
	if err != nil {
		return serializer.Response{
			Code:  50003,
			Msg:   "Redis 连接失败",
			Error: err.Error(),
		}
	}

	if len(result) > 0 {
		return serializer.Response{
			Data: serializer.BuildSensorThresholdByRedis(result),
		}
	}

	db_primary, _ := strconv.Atoi(id)
	var sensorThreshold model.SensorsThreshold
	sensorThreshold.EquipmentID = uint32(db_primary)
	dbErr := model.DB.First(&sensorThreshold, db_primary).Error
	if dbErr != nil {
		return serializer.Response{
			Code:  404,
			Msg:   "数据库无数据,请先设置值",
			Error: dbErr.Error(),
		}
	}
	return serializer.Response{
		Data: serializer.BuildSensorThresholdByDB(sensorThreshold),
	}

}
