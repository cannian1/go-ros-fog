package api

import (
	"go-ros-fog/service"

	"github.com/gin-gonic/gin"
)

func SensorValue(c *gin.Context) {
	service := service.TcpSensorService{}
	res := service.SensorValue(c.Param("id"))
	c.JSON(200, res)
}
