package tcp_model

import (
	"encoding/json"
	"fmt"
	"go-ros-fog/ziface"
	"sync"
)

type DeviceID struct {
	ID uint32 `json:"d_id"`
}

type Sensor struct {
	DeviceID    uint32  `json:"device_id"`
	Temperature float32 `json:"temperature"`
	LightLevel  uint32  `json:"light_intensity"`
	Smog        uint32  `json:"smog"`
	Time        int64   `json:"time"`
}

type Servo struct {
	DeviceID
	State bool    `json:"state"`
	Time  int64   `json:"time"`
	Angle float32 `json:"angle"`
}

type Relay struct {
	DeviceID
	State bool  `json:"state"`
	Time  int64 `json:"time"`
}

type EdgeDevice struct {
	Conn     ziface.IConnection // 连接
	DeviceID uint32             `json:"device_id"`
	Sensor   *Sensor
	Servo    *Servo
	Relay    *Relay
}

var DeviceGen uint32 = 1
var IDLock sync.Mutex

// NewEdgeDevice 初始化边缘设备，分配 id
func NewEdgeDevice(conn ziface.IConnection) *EdgeDevice {
	IDLock.Lock()
	id := DeviceGen
	DeviceGen++
	IDLock.Unlock()

	ed := EdgeDevice{
		Conn:     conn,
		DeviceID: id,
		Sensor:   &Sensor{id, 0, 0, 0, 0},
		Servo:    &Servo{DeviceID{id}, false, 0, 0},
		Relay:    &Relay{DeviceID{id}, false, 0},
	}
	return &ed
}

// SendMsg 向客户端发送消息
func (ed *EdgeDevice) SendMsg(msgID uint32, data interface{}) {
	// 将发送的消息序列化为 json 格式
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}

	if ed.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := ed.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("EdgeDevice SendMsg Error", err)
		return
	}
}

// SyncDid 上线同步设备id
func (ed *EdgeDevice) SyncDid() {
	ed.SendMsg(1, ed.DeviceID)
}

// Marshal2Sensor 将 JSON 字节流数据序列化为 Sensor 数据
func (ed *EdgeDevice) Marshal2Sensor(data []byte) *Sensor {
	tempSensor := Sensor{}
	err := json.Unmarshal(data, &tempSensor)
	if err != nil {
		fmt.Println(err)
	}
	return &tempSensor
}
