package serializer

import "strconv"

type Amcl struct {
	PosePosePositionX    float32 `json:"pose_pose_position_x"`
	PosePosePositionY    float32 `json:"pose_pose_position_y"`
	PosePoseOrientationZ float32 `json:"pose_pose_orientation_z"`
}

func BuildAmcl(x, y, z string) Amcl {
	tx, _ := strconv.ParseFloat(x, 32)
	ty, _ := strconv.ParseFloat(y, 32)
	tz, _ := strconv.ParseFloat(z, 32)

	return Amcl{
		PosePosePositionX:    float32(tx),
		PosePosePositionY:    float32(ty),
		PosePoseOrientationZ: float32(tz),
	}
}
