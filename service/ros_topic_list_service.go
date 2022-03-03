package service

import (
	"go-ros-fog/serializer"
)

type RosTopicList struct {
}

func (service *RosTopicList) ListSurvivingTopics() serializer.Response {
	// 如果出错，就是连不上rosmaster，直接panic了
	return serializer.Response{
		Data: serializer.BuildTopicListsAlive(),
	}
}
