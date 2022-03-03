package serializer

import "go-ros-fog/ros"

// 构建 rostopic list
func BuildTopicListsAlive() []string {
	topics := ros.GetTopics()
	return topics
}

// 构建 rostopic list pub
func BuildTopicListPub() map[string][]string {
	topics := ros.GetPubTopic()
	return topics
}

// 构建 rostopic list sub
func BuildTopicListSub() map[string][]string {
	topics := ros.GetSubTopic()
	return topics
}

// 构建历史 rostopic list
func BuildTopicListHistory() []string {
	topics := ros.GetHistoryTopics()
	return topics
}
