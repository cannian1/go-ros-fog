package model

type OOB struct {
	Date                    int `json:"date" gorm:"primaryKey;autoIncrement:false"`
	TemperatureAbnormalTime int    `json:"temperature_at"`
	LightLevelAbnormalTime  int    `json:"lightlevel_at"`
	SmogAbnormalTime        int    `json:"smog_at"`
}
