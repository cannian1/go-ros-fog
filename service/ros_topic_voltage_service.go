package service

import (
	"go-ros-fog/cache"
	"go-ros-fog/serializer"
)

type RosTopicVoltage struct {
}

func (service *RosTopicVoltage) GetVoltage() serializer.Response {
	result, err := cache.RedisClient.Get(cache.RosTopicPowerVoltage).Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Data: result,
	}
}
