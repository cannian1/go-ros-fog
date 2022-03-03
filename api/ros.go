package api

import (
	"go-ros-fog/service"

	"github.com/gin-gonic/gin"
)

// RosTopicList 列出当前活跃的 ros 话题
func RosTopicList(c *gin.Context) {
	service := service.RosTopicList{}
	res := service.ListSurvivingTopics()
	c.JSON(200, res)
}

// RosTopicPubs 列出当前发布的 ros 话题
func RosTopicPubs(c *gin.Context) {
	service := service.RosTopicList{}
	res := service.ListPubTopics()
	c.JSON(200, res)
}

// RosTopicSubs 列出当前订阅的 ros 话题
func RosTopicSubs(c *gin.Context) {
	service := service.RosTopicList{}
	res := service.ListSubTopics()
	c.JSON(200, res)
}

// RosTopicListHistory 列出历史注册过的 ros 话题
func RosTopicListHistory(c *gin.Context) {
	service := service.RosTopicList{}
	res := service.ListHistoryTopics()
	c.JSON(200, res)
}

// RosMachines 列出当前所有连接到ROS Master的机器地址
func RosMachines(c *gin.Context) {

}

// RosNodeList 列出所有ros节点
func RosNodeList(c *gin.Context) {

}

// RosServiceList 列出所有ros服务
func RosServiceList(c *gin.Context) {

}
