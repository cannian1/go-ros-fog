package ros

import (
	"fmt"
	"go-ros-fog/cache"
	"go-ros-fog/ros_model"
	"time"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/nav_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type BusinessNode struct {
	subOdomTopic         *goroslib.Subscriber // 订阅 Odom 话题
	subAmclPoseTopic     *goroslib.Subscriber // 订阅 amcl 话题
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

	subOdomTopic, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     rosNode,
		Topic:    "odom",
		Callback: SubOdomCallBack,
	})
	if err != nil {
		panic(err)
	}

	bn.subAmclPoseTopic = subAmclPose
	bn.subChatterTopic = subChatter
	bn.sub2DNavGoalTopic = subGoal
	bn.subPowerVoltageTopic = subPowerVoltage
	bn.subOdomTopic = subOdomTopic
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
				DeviceId:    1,
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
	// redis 设置 key value 过期时间
	err := cache.RedisClient.Set(cache.RosTopicCharrter, msg.Data, 5*time.Second).Err()
	if err != nil {
		panic("[redis error]" + err.Error())
	}
}

// SubOdomCallBack Odom回调函数,设置 X、Y轴线速度
func SubOdomCallBack(msg *nav_msgs.Odometry) {
	// redis 设置 key value 过期时间
	if err := cache.RedisClient.Set(cache.RosTopicOdom+":twist_x", msg.Twist.Twist.Linear.X, 1*time.Second).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}

	if err := cache.RedisClient.Set(cache.RosTopicOdom+":twist_y", msg.Twist.Twist.Linear.Y, 1*time.Second).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}

}

// SubAmclPoseCallBack amcl_pose 回调函数
func SubAmclPoseCallBack(msg *geometry_msgs.PoseWithCovarianceStamped) {
	if err := cache.RedisClient.Set(cache.RosTopicAmclPose+":pose_pose_position_x", msg.Pose.Pose.Position.X, 10*time.Minute).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}

	if err := cache.RedisClient.Set(cache.RosTopicAmclPose+":pose_pose_position_y", msg.Pose.Pose.Position.Y, 10*time.Minute).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}

	if err := cache.RedisClient.Set(cache.RosTopicAmclPose+":pose_pose_orientation_z", msg.Pose.Pose.Orientation.Z, 10*time.Minute).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}
}

// Sub2DNavGoalCallBack 导航目标话题回调
func Sub2DNavGoalCallBack(msg *geometry_msgs.PoseStamped) {
	// Setting goal: Frame:map, Position(-2.967, -0.971, 0.000), Orientation(0.000, 0.000, 0.373, 0.928) = Angle: 0.764

	if err := cache.RedisClient.Set(cache.RosTopicMoveBaseGoal+":pose_position_x", msg.Pose.Position.X, 30*time.Minute).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}

	if err := cache.RedisClient.Set(cache.RosTopicMoveBaseGoal+":pose_position_y", msg.Pose.Position.Y, 30*time.Minute).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}

	if err := cache.RedisClient.Set(cache.RosTopicMoveBaseGoal+":pose_orientation_z", msg.Pose.Orientation.Z, 30*time.Minute).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}

	if err := cache.RedisClient.Set(cache.RosTopicMoveBaseGoal+":pose_orientation_w", msg.Pose.Orientation.W, 30*time.Minute).Err(); err != nil {
		panic("[redis error]" + err.Error())
	}
}

// SubPowerVoltageCallBack 电压回调
func SubPowerVoltageCallBack(msg *std_msgs.Float32) {
	// redis 设置 key value 过期时间
	err := cache.RedisClient.Set(cache.RosTopicCharrter, msg.Data, 1*time.Minute).Err()
	if err != nil {
		panic("[redis error]" + err.Error())
	}
}

// CloseSub 回收资源
func (bn *BusinessNode) CloseSub() {
	bn.subAmclPoseTopic.Close()
	bn.subChatterTopic.Close()
	bn.subGoalTopic.Close()
	bn.subPowerVoltageTopic.Close()
	bn.subOdomTopic.Close()
}

func (bn *BusinessNode) ClosePub() {
	bn.pubSensorTopic.Close()
}
