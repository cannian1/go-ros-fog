package serializer

import "go-ros-fog/cache"

func BuildSensor() map[string]string {
	res, _ := cache.RedisClient.HGetAll(cache.SensorValue).Result()
	return res
}



// func BuildSensors() []map[string]string{
	
// }