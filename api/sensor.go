package api

import (
	"encoding/json"
	"fmt"
	"go-ros-fog/service"
	"log"

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
	//recv := make(map[string]interface{})

	for i := 0; i < n; i++ {
		if buf[i] == ' ' || buf[i] == '\x00' {
			buf = append(buf[:i], buf[(i+1):]...)
			n -= 1
		}
	}

	t := buf[0:n]

	//
	// if err := json.Unmarshal(t, &recv); err == nil {
	// 	t := make(map[string]interface{})
	// 	t["equipment_id"] = c.Param("id")
	// 	t["temperature"] = recv["temperature"]
	// 	t["light_intensity"] = recv["light_intensity"]
	// 	t["smog"] = recv["smog"]
	// 	res := service.SetSensorThreshold(c.Param("id"), t)
	// 	c.JSON(200, res)
	// } else {
	// 	log.Println("接收到异常body")
	// 	c.JSON(200, ErrorResponse(err))
	// }
	res := service.SetSensorThreshold(c.Param("id"), t)
	c.JSON(200, res)
}

func SensorThreshold_old(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	recv := make(map[string]interface{})
	fmt.Println(string(buf[0:n]))
	for i := 0; i < n; i++ {
		if buf[i] == ' ' || buf[i] == '\x00' {
			buf = append(buf[:i], buf[(i+1):]...)
			n -= 1
		}
	}

	t := buf[0:n]
	err := json.Unmarshal(t, &recv)
	if err != nil {
		log.Println("接收到异常body")
	}
	c.JSON(200, recv)

}
