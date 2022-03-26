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
	light, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.No2AbnormalTime).Result()
	smog, _ := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.CoAbnormalTime).Result()
	t, _ := strconv.Atoi(temp)
	l, _ := strconv.Atoi(light)
	s, _ := strconv.Atoi(smog)

	key, _ := strconv.Atoi(cache.GetDateKey())
	oob := model.OOB{
		Date:                    key,
		TemperatureAbnormalTime: t,
		No2AbnormalTime:         l,
		CoAbnormalTime:          s,
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

func (service *TcpSensorAbnormalTimeService) GetOutOfBorderNow() serializer.Response {
	temp, err := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.TemperatureAbnormalTime).Result()
	if err != nil {
		return serializer.Response{
			Code:  200,
			Msg:   "Redis 暂无越界数据",
			Error: err.Error(),
		}
	}

	light, err := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.No2AbnormalTime).Result()
	if err != nil {
		return serializer.Response{
			Code:  200,
			Msg:   "Redis 暂无越界数据",
			Error: err.Error(),
		}
	}

	smog, err := cache.RedisClient.Get(cache.GetDateKey() + ":" + cache.CoAbnormalTime).Result()
	if err != nil {
		return serializer.Response{
			Code:  200,
			Msg:   "Redis 暂无越界数据",
			Error: err.Error(),
		}
	}

	if len(temp) > 0 || len(light) > 0 || len(smog) > 0 {
		return serializer.Response{
			Data: serializer.BuildOutOfBorderNow(temp, light, smog),
		}
	}
	return serializer.Response{
		Code: 200,
		Msg:  "",
	}
}

func (service *TcpSensorAbnormalTimeService) GetOutOfBorderLast7Days() serializer.Response {
	var oobs []model.OOB

	err := model.DB.Where("date > ? - 7 ", cache.GetDateKey()).Find(&oobs).Error
	if err != nil {
		return serializer.Response{
			Code:  50000,
			Msg:   "数据库连接错误",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Data: serializer.BuildOutOfBorderLast7Days(oobs),
	}
}
