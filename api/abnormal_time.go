package api

import (
	"go-ros-fog/service"

	"github.com/gin-gonic/gin"
)

func DelOutOfBorder(c *gin.Context) {
	service := service.TcpSensorAbnormalTimeService{}
	res := service.DelHistoryData(c.Param("id"))
	c.JSON(200, res)
}

func DataFlush2DB(c *gin.Context) {
	service := service.TcpSensorAbnormalTimeService{}
	res := service.Flush2DB()
	c.JSON(200, res)
}
