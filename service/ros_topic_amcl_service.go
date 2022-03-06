package service

import (
	"go-ros-fog/cache"
	"go-ros-fog/serializer"
)

type RosTopicAmcl struct {
}

func (service *RosTopicAmcl) GetAmcl() serializer.Response {
	positionX, err := cache.RedisClient.Get(cache.RosTopicAmclPose + ":pose_pose_position_x").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
			Msg: "小车在10分钟之内没有移动超过0.2米或旋转超过30度",
		}
	}

	positionY, err := cache.RedisClient.Get(cache.RosTopicAmclPose + ":pose_pose_position_y").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
			Msg: "小车在10分钟之内没有移动超过0.2米或旋转超过30度",
		}
	}

	oriZ, err := cache.RedisClient.Get(cache.RosTopicAmclPose + ":pose_pose_orientation_z").Result()
	if err != nil {
		return serializer.Response{
			Code:  serializer.CodeRedisDBErr,
			Error: err.Error(),
			Msg: "小车在10分钟之内没有移动超过0.2米或旋转超过30度",
		}
	}

	return serializer.Response{
		Data: serializer.BuildAmcl(positionX, positionY, oriZ),
	}
}
