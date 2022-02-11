package model

import (
	"encoding/json"
	"fmt"
)

type EchoMsg struct {
	EquipmentID uint32 `json:"equipment_id"`
	Msg         string `json:"msg"`
	Time        int64  `json:"time"`
}

// Unmarshal 结构体
func (e *EchoMsg) Unmarshal(data []byte) {
	err := json.Unmarshal([]byte(data), &e)
	if err != nil {
		fmt.Println("反序列化错误")
	}
}

// Marshal 结构体
func (e *EchoMsg) Marshal() (data []byte) {
	data, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}
	return data
}
