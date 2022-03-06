package api

import (
	"go-ros-fog/service"

	"github.com/gin-gonic/gin"
)

func RosTopicRecvChatter(c *gin.Context) {
	service := service.RosTopicChatter{}
	res := service.GetChatter()
	c.JSON(200, res)
}

func RosTopicRecvVoltage(c *gin.Context) {
	service := service.RosTopicVoltage{}
	res := service.GetVoltage()
	c.JSON(200, res)
}

func RosTopicRecvOdom(c *gin.Context) {
	service := service.RosTopicOdom{}
	res := service.GetOdom()
	c.JSON(200, res)
}

func RosTopicRecvAmcl(c *gin.Context) {
	service := service.RosTopicAmcl{}
	res := service.GetAmcl()
	c.JSON(200, res)
}

func RosTopicRecvGoal(c *gin.Context){
	service := service.RosTopicGoal{}
	res := service.GetGoal()
	c.JSON(200, res)
}