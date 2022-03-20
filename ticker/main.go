package ticker

import "go-ros-fog/tcp_core"

func HandleAllTicker() {
	// 每天0点数据持久化到磁盘，并清理缓存
	SetTickerEveryDay(tcp_core.Save2DB)
}
