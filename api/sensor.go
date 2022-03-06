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

func SetSensorThreshold(c *gin.Context) {
	service := service.TcpSensorService{}

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	for i := 0; i < n; i++ {
		if buf[i] == ' ' || buf[i] == '\x00' {
			buf = append(buf[:i], buf[(i+1):]...)
			n -= 1
		}
	}

	t := buf[0:n]

	res := service.SetSensorThreshold(c.Param("id"), t)
	c.JSON(200, res)
}

func GetSensorThreshold(c *gin.Context) {
	service := service.TcpSensorService{}
	res := service.GetSensorThreshold(c.Param("id"))
	c.JSON(200, res)
}