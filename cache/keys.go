package cache

import (
	"fmt"
	"strconv"
)

const (
	// tcp 连接数
	TCPConnCount = "tcp_conn:count"
	// tcp 连接服务类型
	TCPConnServiceType = "tcp_conn:service_type"
	// 每小时排行
	TemperatureHourlyRankKey = "temp_rank:hourly"
)

// TemperatureSensorKey 传感器 TCP 传输来的温度
// sensor:temperature:1 -> 21
// sensor:temperature:2 -> 30
func TemperatureSensorKey(id uint) string {
	return fmt.Sprintf("sensor:temperature:%s", strconv.Itoa(int(id)))
}
