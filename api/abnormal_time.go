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
