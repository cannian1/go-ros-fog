package ticker

import "go-ros-fog/tcp_core"

func HandleAllTicker() {
	// 每天0点数据持久化到磁盘，并清理缓存
	defer SetTickerEveryDay(tcp_core.Save2DBAndPurge)
	// 每 hour 数据持久化到磁盘
	SetTickerEveryHour(tcp_core.Save2DB)
}
