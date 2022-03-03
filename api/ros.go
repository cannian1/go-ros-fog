package api

import (
	"go-ros-fog/service"

	"github.com/gin-gonic/gin"
)

// TODO:
// RosTopicList 列出当前活跃的 ros 话题
func RosTopicList(c *gin.Context) {
	service := service.RosTopicList{}
	res := service.ListSurvivingTopics()
	c.JSON(200, res)
}
