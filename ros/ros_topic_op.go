package ros

import (
	"fmt"
	"go-ros-fog/ros_model"
	"time"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type BusinessNode struct {
	subAmclPoseTopic     *goroslib.Subscriber //
	subGoalTopic         *goroslib.Subscriber // 订阅目标话题
	subChatterTopic      *goroslib.Subscriber
	sub2DNavGoalTopic    *goroslib.Subscriber // move_base_simple/goal
	subPowerVoltageTopic *goroslib.Subscriber // 订阅电压话题
	pubSensorTopic       *goroslib.Publisher  // 发布自定义传感器话题
}

func (bn *BusinessNode) InitSubscriber() {
	// 订阅当前位姿话题
	subAmclPose, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     rosNode,
		Topic:    "amcl_pose",
		Callback: SubAmclPoseCallBack,
	})
	if err != nil {
		panic(err)
	}

	// 订阅 chatter 话题
	subChatter, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     rosNode,
		Topic:    "chatter",
		Callback: SubChatterCallBack,
	})
	if err != nil {
		panic(err)
	}

	// 订阅目标话题
	subGoal, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     rosNode,
		Topic:    "move_base_simple/goal",
		Callback: Sub2DNavGoalCallBack,
	})
	if err != nil {
		subGoal.Close()
		panic(err)
	}

	subPowerVoltage, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     rosNode,
		Topic:    "PowerVoltage",
		Callback: SubPowerVoltageCallBack,
	})
	if err != nil {
		panic(err)
	}

	bn.subAmclPoseTopic = subAmclPose
	bn.subChatterTopic = subChatter
	bn.sub2DNavGoalTopic = subGoal
	bn.subPowerVoltageTopic = subPowerVoltage
}

// InitPublisher 初始化发布者话题
func (bn *BusinessNode) InitPublisher() {
	// create a publisher
	pub, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  rosNode,
		Topic: "sensors_msg",
		Msg:   &ros_model.Sensors{},
	})
	if err != nil {
		pub.Close()
		panic(err)
	}
	bn.pubSensorTopic = pub

	r := rosNode.TimeRate(1 * time.Second)

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)

	for {
		select {
		// publish a message every second
		case <-r.SleepChan():
			msg := &ros_model.Sensors{
				DeviceId:    20,
				Temperature: 30,
				LightLevel:  26,
				Smog:        47,
				Time:        time.Now().UnixNano(),
			}
			fmt.Printf("Outgoing: %+v\n", msg)
			pub.Write(msg)

			// handle CTRL-C
			//case <-c:
			//	return
		}
	}
}

// SubChatterCallBack Chatter 回调函数
func SubChatterCallBack(msg *std_msgs.String) {
	fmt.Printf("Incoming: %+v\n", msg.Data)
}

// SubAmclPoseCallBack amcl_pose 回调函数
func SubAmclPoseCallBack(msg *geometry_msgs.PoseWithCovarianceStamped) {
	fmt.Printf("amcl %+v\n", msg)
}

// Sub2DNavGoalCallBack 导航目标话题回调
func Sub2DNavGoalCallBack(msg *geometry_msgs.PoseStamped) {
	fmt.Printf("goal: %+v\n", msg)
}

// SubPowerVoltageCallBack 电压回调
func SubPowerVoltageCallBack(msg *std_msgs.Float32) {
	fmt.Printf("goal: %+v\n", msg)
}

// CloseSub 回收资源
func (bn *BusinessNode) CloseSub() {
	bn.subAmclPoseTopic.Close()
	bn.subChatterTopic.Close()
	bn.subGoalTopic.Close()
	bn.subPowerVoltageTopic.Close()
}

func (bn *BusinessNode) ClosePub() {
	bn.pubSensorTopic.Close()
}
