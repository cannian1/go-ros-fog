package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 测试 datapack 拆包 封包的单元测试
// 执行单元测试的时候，记得把 utils/globalobj.go 里的 GlobalObject.Reload() 方法注释掉
func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 创建 socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	// 创建一个 go ，负责从客户端处理业务
	go func() {
		// 从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				// -------拆包的过程--------
				// 定义一个拆包的对象 dp
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					// io.ReadFull 读 len(headData)的数据
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						break
					}
					if msgHead.GetMsgLen() > 0 {
						// msg 是有数据的，需要进行第二次读取
						// 第二次从conn读，根据head中的dataLen，再读取data内存
						msg := msgHead.(*Message) // interface 类型的数据，断言
						msg.data = make([]byte, msg.GetMsgLen())

						// 根据 dataLen 的长度再次从 io流中读取
						_, err := io.ReadFull(conn, msg.data)
						{
							if err != nil {
								fmt.Println("server unpack data err", err)
								return
							}
							// 完整的消息已经读取完毕
							fmt.Println("------->Recv MsgID:", msg.id, "dataLen = ", msg.dataLen, "data = ", string(msg.data))
						}
					}
				}

			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	// 创建一个封包对象 dp
	dp := NewDataPack()
	// 模拟“粘包”过程，封装两个msg一起发送
	msg1 := &Message{
		id:      1,
		dataLen: 4,
		data:    []byte{'z', 'i', 'n', 'x'},
	}

	msg2 := &Message{
		id:      2,
		dataLen: 7,
		data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}

	// 封装第一个msg1包
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}

	// 封装第二个msg2包
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error", err)
		return
	}
	// 将两个包 粘 在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
