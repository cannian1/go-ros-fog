package service

import (
	"go-ros-fog/serializer"
)

// 算了这就不写注释了
type RosTopicList struct {
}

func (service *RosTopicList) ListSurvivingTopics() serializer.Response {
	// 如果出错，就是连不上rosmaster，直接panic了
	return serializer.Response{
		Data: serializer.BuildTopicListsAlive(),
	}
}

func (service *RosTopicList) ListSubTopics() serializer.Response {
	return serializer.Response{
		Data: serializer.BuildTopicListSub(),
	}
}

func (service *RosTopicList) ListPubTopics() serializer.Response {
	return serializer.Response{
		Data: serializer.BuildTopicListPub(),
	}
}

func (service *RosTopicList) ListHistoryTopics() serializer.Response {
	return serializer.Response{
		Data: serializer.BuildTopicListHistory(),
	}
}
