package model

type OOB struct {
	Date                    int `json:"date" gorm:"primaryKey;autoIncrement:false"`
	TemperatureAbnormalTime int `json:"temperature_at"`
	No2AbnormalTime         int `json:"no2_at"`
	CoAbnormalTime          int `json:"co_at"`
}
