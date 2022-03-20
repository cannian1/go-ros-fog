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

	if err := cache.RedisClient.Del(cache.GetDateKey() + ":" + cache.LightAbnormalTime).Err(); err != nil {
		panic(err.Error() + "清除光强每日异常时长缓存失败")
	}

	if err := cache.RedisClient.Del(cache.GetDateKey() + ":" + cache.SmogAbnormalTime).Err(); err != nil {
		panic(err.Error() + "清除烟雾每日异常时长缓存失败")
	}

	if err := cache.RedisClient.Del(cache.GetDateKey() + ":" + cache.AbnormalTime).Err(); err != nil {
		panic(err.Error() + "每日异常时长缓存失败")
	}
}

func Save2DB() {
	defer ClearDailyCache()

	temp, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.TemperatureAbnormalTime).Result()
	light, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.LightAbnormalTime).Result()
	smog, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.SmogAbnormalTime).Result()
	t, _ := strconv.Atoi(temp)
	l, _ := strconv.Atoi(light)
	s, _ := strconv.Atoi(smog)

	key, _ := strconv.Atoi(cache.GetDateKey())
	oob := model.OOB{
		Date:                    key,
		TemperatureAbnormalTime: t,
		LightLevelAbnormalTime:  l,
		SmogAbnormalTime:        s,
	}

	err := model.DB.Create(&oob).Error
	if err != nil {
		panic("[fatal err]" + err.Error())
	}

}
