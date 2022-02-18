package tcp_model

import (
	"encoding/json"
	"fmt"
)

// TODO: 更改返回类型
type ResRelay struct {
	Temperature bool  `json:"res_temperature"`
	LightLevel  bool  `json:"res_light_intensity"`
	Smog        bool  `json:"res_smog"`
	Time        int64 `json:"set_time"`
}

func (rr *ResRelay) Marshal() []byte {
	res, err := json.Marshal(&rr)
	if err != nil {
		fmt.Println(err)
	}
	return res
}
