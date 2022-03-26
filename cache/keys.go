package cache

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// tcp 连接数
	TCPConnCount = "tcp_conn:count"
	// tcp 连接服务类型
	TCPConnServiceType = "tcp_conn:service_type"
	// 每小时排行
	TemperatureHourlyRankKey = "temp_rank:hourly"
	// 传感器采集数据值
	SensorValue = "tcp_conn:sensor_value"
	// 传感器阈值
	SensorThreshold = "sensor_threshold"
	// ros话题 chatter
	RosTopicCharrter = "rostopic:chatter"
	// ros话题 amcl_pose
	RosTopicAmclPose = "rostopic:amcl_pose"
	// ros话题 PowerVoltage
	RosTopicPowerVoltage = "rostopic:PowerVoltage"
	// ros话题 move_base_simple/goal
	RosTopicMoveBaseGoal = "rostopic:move_base_simple/goal"
	// ros话题 Odom
	RosTopicOdom = "rostopic:odom"

	AbnormalTime = "abnormal_time"

	TemperatureAbnormalTime = "temperature_abnormal_time"

	No2AbnormalTime = "no2_abnormal_time"

	CoAbnormalTime = "co_abnormal_time"
)

// TemperatureSensorKey 传感器 TCP 传输来的温度
// sensor:temperature:1 -> 21
// sensor:temperature:2 -> 30
func TemperatureSensorKey(id uint) string {
	return fmt.Sprintf("sensor:temperature:%s", strconv.Itoa(int(id)))
}

// GetDateKey 返回当前日期作为 Key
func GetDateKey() string {
	return time.Now().Format("20060102")
}
