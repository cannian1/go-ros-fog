package model

import (
	"encoding/json"
	"fmt"
)

type Sensors struct {
	//gorm.Model
	Temperature    float32 `json:"temperature"`
	LightIntensity int32   `json:"light_intensity"`
	Smog           int32   `json:"smog"`
	Time           int64   `json:"time"`
}

// Unmarshal 反序列化json to Sensors 结构体
func (s *Sensors) Unmarshal(data []byte) {
	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		fmt.Println("反序列化错误")
	}
}
