package service

import (
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
			Code:  404,
			Msg:   "数据库无数据",
			Error: dbErr.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Msg:  "ok",
	}
}
