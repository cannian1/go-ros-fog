package ros

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aler9/goroslib"
)

const ROS_MASTER_ADDRESS = "192.168.0.104:11311"

// ROSNode 声明一个 ROS 节点句柄
var rosNode *goroslib.Node

func init() {
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "goroslib_0",
		MasterAddress: ROS_MASTER_ADDRESS,
	})

	// 连接到 ROS 节点失败
	if err != nil {
		defer fmt.Println("连接到 ROS 节点失败!请检查ip和端口，以及对应主机ros master状况")
		n.Close()
		panic(err)
	}
	rosNode = n
}

// CloseNode 关闭该程序创建的节点
func CloseNode() {
	fmt.Println("连接到 ROS 节点失败!请检查ip和端口，以及对应主机ros master状况")
	rosNode.Close()
}

// GetNodes 获取所有连接到 Ros Master 的节点
func GetNodes() []string {
	nodes, err := rosNode.MasterGetNodes()
	if err != nil {
		CloseNode()
		panic(err)
	}

	nodeList := make([]string, 0, 2)
	for i := range nodes {
		nodeList = append(nodeList, i)
	}
	return nodeList
}

// GetMachines 获取所有连接到 ROS Master 的机器
func GetMachines() []string {
	machines, err := rosNode.MasterGetMachines()
	if err != nil {
		CloseNode()
		panic(err)
	}

	machineList := make([]string, 0, 2)
	for i := range machines {
		machineList = append(machineList, i)
	}
	return machineList
}

// GetHistoryTopics 获取历史在ros master中注册的话题
func GetHistoryTopics() []string {
	topics, err := rosNode.MasterGetTopics()
	if err != nil {
		CloseNode()
		panic(err)
	}

	topicList := make([]string, 0, 2)
	for i := range topics {
		topicList = append(topicList, i)
	}
	return topicList
}

/*
	type InfoTopic struct {
	Type        string
	Publishers  map[string]struct{}
	Subscribers map[string]struct{}
}

		ret[entry.Name] = &InfoTopic{
			Type:        entry.Type,
			Publishers:  make(map[string]struct{}),
			Subscribers: make(map[string]struct{}),
		}
*/

func GetPubTopic() map[string][]string {
	topics, err := rosNode.MasterGetTopics()
	if err != nil {
		CloseNode()
		panic(err)
	}
	// key 是话题名，value是被哪个节点发布
	topicPub := make(map[string][]string, 0)
	for i := range topics {
		for j := range topics[i].Publishers {
			topicPub[i] = append(topicPub[i], j)
		}
	}
	return topicPub
}

// GetSubTopic 根据话题找到哪个节点订阅的
func GetSubTopic() map[string][]string {
	topics, err := rosNode.MasterGetTopics()
	if err != nil {
		CloseNode()
		panic(err)
	}
	// key 是话题名，value是被哪个节点订阅
	topicSub := make(map[string][]string, 0)
	for i := range topics {
		for j := range topics[i].Subscribers {
			topicSub[i] = append(topicSub[i], j)
		}
	}
	return topicSub
}

// GetTopics 获取所有连接到 ROS Master 的正在订阅或发布的话题
func GetTopics() []string {
	topics, err := rosNode.MasterGetTopics()
	if err != nil {
		CloseNode()
		panic(err)
	}

	topicList := make([]string, 0, 2)

	// 遍历所有话题，如果有发布或订阅信息，就添加到切片中
	for i := range topics {
		if _, ok := GetPubTopic()[i]; ok {
			topicList = append(topicList, i)
			continue
		}

		if _, ok := GetSubTopic()[i]; ok {
			topicList = append(topicList, i)
			continue
		}
	}

	return topicList
}

// GetServices 获取当前所有连接 ros master 的服务
func GetServices() []string {
	services, err := rosNode.MasterGetServices()
	if err != nil {
		CloseNode()
		panic(err)
	}

	serviceList := make([]string, 0, 2)
	for i := range services {
		serviceList = append(serviceList, i)
	}
	return serviceList
}

func JSONOutputStdout(v interface{}) {
	res, err := json.Marshal(&v)
	if err != nil {
		panic("一定是空指针了")
	}
	var out bytes.Buffer
	json.Indent(&out, res, "", "\t")
	out.WriteTo(os.Stdout)
}
