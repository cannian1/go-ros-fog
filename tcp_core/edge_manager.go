package tcp_core

import (
	"go-ros-fog/tcp_model"
	"sync"
)

type EdgeManager struct {
	edgeDevicesMap map[uint32]*tcp_model.EdgeDevice
	eLock          sync.RWMutex
}

// EDMgrObj 对外的唯一设备管理句柄
var EDMgrObj *EdgeManager

func init() {
	EDMgrObj = &EdgeManager{
		edgeDevicesMap: make(map[uint32]*tcp_model.EdgeDevice),
	}
}

// AddEdgeDevice 添加终端设备到管理器
func (em *EdgeManager) AddEdgeDevice(e *tcp_model.EdgeDevice) {
	em.eLock.Lock()
	em.edgeDevicesMap[e.DeviceID] = e
	em.eLock.Unlock()
}

// Remove 从管理器移除终端设备
func (em *EdgeManager) Remove(did uint32) {
	em.eLock.Lock()
	delete(em.edgeDevicesMap, did)
	em.eLock.Unlock()
}

// GetEdgeDeviceByDId 通过设备ID查询 设备对象
func (em *EdgeManager) GetEdgeDeviceByDId(did uint32) *tcp_model.EdgeDevice {
	em.eLock.Lock()
	defer em.eLock.Unlock()
	// TODO 判断id是否存在
	return em.edgeDevicesMap[did]
}

// GetAllEdgeDevice 获取所有边缘设备
func (em *EdgeManager) GetAllEdgeDevice() []*tcp_model.EdgeDevice {
	em.eLock.Lock()
	defer em.eLock.Unlock()

	devices := make([]*tcp_model.EdgeDevice, 0, len(em.edgeDevicesMap))
	for _, v := range em.edgeDevicesMap {
		devices = append(devices, v)
	}
	return devices
}
