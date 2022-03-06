package service

import (
	"go-ros-fog/cache"
	"go-ros-fog/serializer"
)

type RosTopicOdom struct {
}

func (service *RosTopicOdom) GetOdom() serializer.Response {
	twistX, err := cache.RedisClient.Get(cache.RosTopicOdom + ":twist_x").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
		}
	}

	twistY, err := cache.RedisClient.Get(cache.RosTopicOdom + ":twist_y").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Data: serializer.BuildOdom(twistX, twistY),
	}
}
