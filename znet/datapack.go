package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"go-ros-fog/util"
	"go-ros-fog/ziface"
)

/*
	head		 |      	body    	| head | body
	dataLen | id |			data		| ......
*/

// DataPack 封包、拆包的模块
type DataPack struct {
}

// NewDataPack 实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// dataLen uint32(4字节) + id uint32(4字节）
	return 8
}

// Pack 封包方法
// |dataLen|msgId|data|
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{}) // 空的字节切片
	// 将 dataLen 写进 dataBuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将 MsgId 写进 dataBuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将 data数据 写进 dataBuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// UnPack 拆包方法
// 将包的Head信息读出来，再根据head信息里的data长度，再进行读
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压 head 信息,得到 dataLen 和 MsgID,填充到 msg 中
	msg := &Message{}

	// 读 dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}
	// 读 MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}

	// 判断 dataLen 是否已经超出了允许的最大包长度
	if util.GlobalObject.MaxPackageSize > 0 && msg.dataLen > util.GlobalObject.MaxPackageSize {
		return nil, errors.New("too long data received")
	}
	return msg, nil
}
