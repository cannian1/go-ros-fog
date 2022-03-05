package serializer

import (
	"go-ros-fog/model"
	"strconv"
)

type SensorThreshold struct {
	EquipmentID    uint32  `json:"equipment_id"`
	Temperature    float32 `json:"temperature"`
	LightIntensity uint32  `json:"light_intensity"`
	Smog           uint32  `json:"smog"`
}

// 序列化redis查出的map
func BuildSensorThresholdByRedis(result map[string]string) SensorThreshold {
	equipmentID, _ := strconv.Atoi(result["equipment_id"])
	temperature, _ := strconv.ParseFloat(result["temperature"], 32)
	light_intensity, _ := strconv.ParseUint(result["light_intensity"], 10, 32)
	smog, _ := strconv.ParseUint(result["smog"], 10, 32)
	return SensorThreshold{
		EquipmentID:    uint32(equipmentID),
		Temperature:    float32(temperature),
		LightIntensity: uint32(light_intensity),
		Smog:           uint32(smog),
	}
}

// 序列化从数据库查出的数据
func BuildSensorThresholdByDB(item model.SensorsThreshold) SensorThreshold {
	return SensorThreshold{
		EquipmentID:    item.EquipmentID,
		Temperature:    item.Temperature,
		LightIntensity: item.LightIntensity,
		Smog:           item.Smog,
	}
}
