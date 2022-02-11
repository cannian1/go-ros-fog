package model

import (
	"encoding/json"
	"fmt"
)

type T_Servos struct {
	EquipmentID uint32  `json:"equipment_id"`
	State       bool    `json:"state"`
	Time        int64   `json:"time"`
	Angle       float32 `json:"angle"`
}

// Unmarshal 反序列化json to T_Servo 结构体
func (s *T_Servos) Unmarshal(data []byte) {
	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		fmt.Println("反序列化错误")
	}
}
