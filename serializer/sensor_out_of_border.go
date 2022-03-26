package serializer

import (
	"go-ros-fog/cache"
	"go-ros-fog/model"
	"go-ros-fog/tcp_core"
	"strconv"
)

type SensorsOutOfBorder struct {
	Date                    int `json:"date" gorm:"primaryKey;autoIncrement:false"`
	TemperatureAbnormalTime int `json:"temperature_at"`
	No2AbnormalTime         int `json:"no2_at"`
	CoAbnormalTime          int `json:"co_at"`
}

func BuildOutOfBorderNow(temp, light, smog string) SensorsOutOfBorder {
	d, _ := strconv.Atoi(cache.GetDateKey())
	t, _ := strconv.Atoi(temp)
	l, _ := strconv.Atoi(light)
	s, _ := strconv.Atoi(smog)

	return SensorsOutOfBorder{
		Date:                    d,
		TemperatureAbnormalTime: t,
		No2AbnormalTime:         l,
		CoAbnormalTime:          s,
	}

}

// BuildSensorsOutOfBorder7 序列化
func BuildSensorsOutOfBorder7(item model.OOB) SensorsOutOfBorder {
	return SensorsOutOfBorder{
		Date:                    item.Date,
		TemperatureAbnormalTime: item.TemperatureAbnormalTime,
		No2AbnormalTime:         item.No2AbnormalTime,
		CoAbnormalTime:          item.CoAbnormalTime,
	}
}

// BuildOutOfBorderLast7Days 序列化列表
func BuildOutOfBorderLast7Days(items []model.OOB) (soobs []SensorsOutOfBorder) {
	tcp_core.Save2DB()
	for _, item := range items {
		soob := BuildSensorsOutOfBorder7(item)
		soobs = append(soobs, soob)
	}
	return soobs
}
