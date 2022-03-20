package ticker

import "time"

// SetTickerEveryDay 每天 0点执行函数
func SetTickerEveryDay(f func()) {
	go func() {
		for {
			// f()
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			f()
		}
	}()
}
