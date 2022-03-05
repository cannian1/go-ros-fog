package model

// 传感器阈值模型
type SensorsThreshold struct {
	EquipmentID    uint32  `json:"equipment_id" gorm:"primaryKey;autoIncrement:false"`
	Temperature    float32 `json:"temperature"`
	LightIntensity uint32  `json:"light_intensity"`
	Smog           uint32  `json:"smog"`
}
