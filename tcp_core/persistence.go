package tcp_core

import (
	"go-ros-fog/cache"
	"go-ros-fog/model"
	"strconv"
)

// StasticalOutOfThreshold 统计越界时长
func StasticalOutOfThreshold(key string) {
	err := cache.RedisClient.Incr(cache.GetDateKey() + ":" + key).Err()
	if err != nil {
		panic(err.Error() + "StasticalOutOfThreshold err")
	}
}

// ClearDailyCache 清除每日异常时长缓存
func ClearDailyCache() {
	if err := cache.RedisClient.Del(cache.GetDateKey() + ":" + cache.TemperatureAbnormalTime).Err(); err != nil {
		panic(err.Error() + "清除温度每日异常时长缓存失败")
	}

	if err := cache.RedisClient.Del(cache.GetDateKey() + ":" + cache.No2AbnormalTime).Err(); err != nil {
		panic(err.Error() + "清除No2每日异常时长缓存失败")
	}

	if err := cache.RedisClient.Del(cache.GetDateKey() + ":" + cache.CoAbnormalTime).Err(); err != nil {
		panic(err.Error() + "清除CO每日异常时长缓存失败")
	}

	if err := cache.RedisClient.Del(cache.GetDateKey() + ":" + cache.AbnormalTime).Err(); err != nil {
		panic(err.Error() + "每日异常时长缓存失败")
	}
}

// Save2DB 从redis中读出越界时长存入数据库中
func Save2DB() {
	temp, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.TemperatureAbnormalTime).Result()
	light, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.No2AbnormalTime).Result()
	smog, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.CoAbnormalTime).Result()
	t, _ := strconv.Atoi(temp)
	l, _ := strconv.Atoi(light)
	s, _ := strconv.Atoi(smog)

	key, _ := strconv.Atoi(cache.GetDateKey())
	oob := model.OOB{
		Date:                    key, //
		TemperatureAbnormalTime: t,
		No2AbnormalTime:         l,
		CoAbnormalTime:          s,
	}

	err := model.DB.Create(&oob).Error
	if err != nil {
		err = model.DB.Save(&oob).Error
		if err != nil {
			panic("[Fatal err]" + err.Error())
		}
	}

}

// Save2DBAndPurge 从redis中读出越界时长存入数据库中 并释放缓存
func Save2DBAndPurge() {
	Save2DB()
	ClearDailyCache()
}
