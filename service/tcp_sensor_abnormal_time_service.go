package service

import (
	"go-ros-fog/cache"
	"go-ros-fog/model"
	"go-ros-fog/serializer"
	"strconv"
)

type TcpSensorAbnormalTimeService struct {
}

func (service *TcpSensorAbnormalTimeService) DelHistoryData(hitherto string) serializer.Response {
	hi, err := strconv.Atoi(hitherto)
	if err != nil {
		return serializer.Response{
			Code: 40001,
			Msg:  "err input",
		}
	}
	var OutOfBorder model.OOB
	dbErr := model.DB.Where("date < ?", hi).Delete(&OutOfBorder).Error
	if dbErr != nil {
		return serializer.Response{
			Code:  serializer.CodeDBError,
			Msg:   "数据库连接异常",
			Error: dbErr.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Msg:  "ok",
	}
}

func (service *TcpSensorAbnormalTimeService) Flush2DB() serializer.Response {
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

	dbErr := model.DB.Create(&oob).Error
	if dbErr != nil {
		dbErr = model.DB.Save(&oob).Error
		if dbErr != nil {
			return serializer.Response{
				Code:  serializer.CodeDBError,
				Msg:   "数据库连接异常",
				Error: dbErr.Error(),
			}
		}
	}

	return serializer.Response{
		Code: 200,
		Msg:  "ok",
	}
}
