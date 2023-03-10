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

		// 用户注册
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
		v1.GET("tcp_sensors/threshold/:id", api.GetSensorThreshold)
		// 获取 Chatter 话题数据
		v1.GET("ros/topic_recv/chatter",api.RosTopicRecvChatter)
		// 获取 电压 话题数据
		v1.GET("ros/topic_recv/voltage",api.RosTopicRecvVoltage)
		// 获取 odom 里程计 话题数据
		v1.GET("ros/topic_recv/odom",api.RosTopicRecvOdom)
		// 获取 amcl_pose 话题数据
		v1.GET("ros/topic_recv/amcl",api.RosTopicRecvAmcl)
		// 获取 move_base_simple/goal 话题数据
		v1.GET("ros/topic_recv/goal",api.RosTopicRecvGoal)

		v1.GET("tcp_sensors/out_of_border",api.GetOutOfBorderNow)
		v1.GET("tcp_sensors/out_of_border/last7days",api.GetOutOfBorderLast7Days)
		// 将缓存中的越界时长统计数据刷到数据库(断电/关机前手动保存数据)
		v1.PUT("tcp_sensors/out_of_border",api.DataFlush2DB)
		// 删除越界统计数据
		v1.DELETE("tcp_sensors/out_of_border/:id",api.DelOutOfBorder)

		// 需要登录保护的
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)

			// auth.DELETE("")
		}
	}
	return r
}
