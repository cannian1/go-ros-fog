package serializer

import (
	"go-ros-fog/model"
	"strconv"
)

type SensorThreshold struct {
	EquipmentID uint32  `json:"equipment_id"`
	Temperature float32 `json:"temperature"`
	No2         float32 `json:"no_2"`
	Co          float32 `json:"co"`
}

// 序列化redis查出的map
func BuildSensorThresholdByRedis(result map[string]string) SensorThreshold {
	equipmentID, _ := strconv.Atoi(result["equipment_id"])
	temperature, _ := strconv.ParseFloat(result["temperature"], 32)
	no2, _ := strconv.ParseFloat(result["no_2"], 32)
	co, _ := strconv.ParseFloat(result["co"], 32)
	return SensorThreshold{
		EquipmentID: uint32(equipmentID),
		Temperature: float32(temperature),
		No2:         float32(no2),
		Co:          float32(co),
	}
}

// 序列化从数据库查出的数据
func BuildSensorThresholdByDB(item model.SensorsThreshold) SensorThreshold {
	return SensorThreshold{
		EquipmentID: item.EquipmentID,
		Temperature: item.Temperature,
		No2:         item.No2,
		Co:          item.Co,
	}
}
