package server

import (
	"go-ros-fog/api"
	"go-ros-fog/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", api.UserLogin)

		// 获取rostopic list
		v1.GET("ros/topic_list/alive", api.RosTopicList)
		v1.GET("ros/topic_list/subscriber", api.RosTopicSubs)
		v1.GET("ros/topic_list/publisher", api.RosTopicPubs)
		v1.GET("ros/topic_list/history", api.RosTopicListHistory)

		// 获取 rosnode list
		v1.GET("ros/node/list", api.RosNodeList)
		// 获取 rosnode machines
		v1.GET("ros/node/machines", api.RosNodeMachines)
		// 获取 rosservice list
		v1.GET("ros/services", api.RosServiceList)

		// 获取传感器的值
		v1.GET("tcp_sensors/get_value/:id", api.SensorValue)
		// 设置传感器阈值
		v1.POST("tcp_sensors/threshold/:id", api.SetSensorThreshold)
		// 获取传感器阈值
		v1.GET("tcp_sensors/threshold/:id",api.GetSensorThreshold)

		// 需要登录保护的
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)
		}
	}
	return r
}
