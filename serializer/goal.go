package serializer

import "strconv"

type Goal struct {
	PosePositionX    float32 `json:"pose_position_x"`
	PosePositionY    float32 `json:"pose_position_y"`
	PoseOrientationZ float32 `json:"pose_orientation_z"`
	PoseOrientationW float32 `json:"pose_orientation_w"`
}

func BuildGoal(x, y, z, w string) Goal{
	tx, _ := strconv.ParseFloat(x, 32)
	ty, _ := strconv.ParseFloat(y, 32)
	tz, _ := strconv.ParseFloat(z, 32)
	tw, _ := strconv.ParseFloat(w, 32)

	return Goal{
		PosePositionX: float32(tx),
		PosePositionY: float32(ty),
		PoseOrientationZ: float32(tz),
		PoseOrientationW: float32(tw),
	}
}
