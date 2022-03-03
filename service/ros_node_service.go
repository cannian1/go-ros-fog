package service

import "go-ros-fog/serializer"

type RosNodeAndServiceEnv struct {
}

func (service *RosNodeAndServiceEnv) ListNodes() serializer.Response {
	// 如果出错，就是连不上rosmaster，直接panic了
	return serializer.Response{
		Data: serializer.BuildNodeList(),
	}
}

func (service *RosNodeAndServiceEnv) ListMachines() serializer.Response {
	return serializer.Response{
		Data: serializer.BuildNodeMachines(),
	}
}

func (service *RosNodeAndServiceEnv) ListServices() serializer.Response {
	return serializer.Response{
		Data: serializer.BuildServiceList(),
	}
}
