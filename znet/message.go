package znet

type Message struct {
	id      uint32 // 消息的ID
	dataLen uint32 // 消息的长度
	data    []byte // 消息的内容
}

// NewMsgPackage 创建一个 Message 消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

// GetMsgId 获取消息的 ID
func (m *Message) GetMsgId() uint32 {
	return m.id
}

// GetMsgLen 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.dataLen
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.data
}

// SetMsgId 设置消息的 ID
func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

// SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.data = data
}

// SetDataLen 设置消息的长度
func (m *Message) SetDataLen(dataLen uint32) {
	m.dataLen = dataLen
}
