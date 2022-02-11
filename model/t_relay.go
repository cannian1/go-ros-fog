package model

import "encoding/json"

type T_Relays struct {
	EquipmentID uint32 `json:"equipment_id"`
	State       bool   `json:"state"`
	Time        int64  `json:"time"`
}

// Unmarshal 反序列化json to T_Relay 结构体
func (tr *T_Relays) Unmarshal(data []byte) {
	json.Unmarshal(data, &tr)
}
