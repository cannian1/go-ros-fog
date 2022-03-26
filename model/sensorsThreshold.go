package model

// 传感器阈值模型
type SensorsThreshold struct {
	EquipmentID uint32  `json:"equipment_id" gorm:"primaryKey;autoIncrement:false"`
	Temperature float32 `json:"temperature"`
	No2         float32 `json:"no_2"`
	Co          float32 `json:"co"`
}
