package service

import (
	"go-ros-fog/cache"
	"go-ros-fog/serializer"
)

type RosTopicGoal struct {
}

func (service *RosTopicGoal) GetGoal() serializer.Response {
	positionX, err := cache.RedisClient.Get(cache.RosTopicMoveBaseGoal + ":pose_position_x").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
			Msg:   "30分钟之内未发布过目标点",
		}
	}

	positionY, err := cache.RedisClient.Get(cache.RosTopicMoveBaseGoal + ":pose_position_y").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
			Msg:   "30分钟之内未发布过目标点",
		}
	}

	oriZ, err := cache.RedisClient.Get(cache.RosTopicMoveBaseGoal + ":pose_orientation_z").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
			Msg:   "30分钟之内未发布过目标点",
		}
	}

	oriW, err := cache.RedisClient.Get(cache.RosTopicMoveBaseGoal + ":pose_orientation_w").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
			Msg:   "30分钟之内未发布过目标点",
		}
	}

	return serializer.Response{
		Data: serializer.BuildGoal(positionX, positionY, oriZ, oriW),
	}
}
