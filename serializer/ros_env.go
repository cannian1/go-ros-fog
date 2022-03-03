package serializer

import "go-ros-fog/ros"

// 构建 rostopic list
func BuildTopicListsAlive() []string{
	topics := ros.GetTopics()
	return topics
}
