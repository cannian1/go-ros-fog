package serializer

import "strconv"

type Odom struct {
	TwistX float32 `json:"twist_x"`
	TwistY float32 `json:"twist_y"`
}

func BuildOdom(x, y string) Odom {
	tx, _ := strconv.ParseFloat(x, 32)
	ty, _ := strconv.ParseFloat(y, 32)

	return Odom{
		TwistX: float32(tx),
		TwistY: float32(ty),
	}
}
