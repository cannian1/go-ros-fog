package ros

import "fmt"

func TopicServer() {
	t := GetNodes()
	fmt.Println(t)
	bussiness := BusinessNode{}
	//time.Sleep(1*time.Second)
	bussiness.InitSubscriber()
	bussiness.InitPublisher()
	defer bussiness.ClosePub()
	defer bussiness.CloseSub()

}
