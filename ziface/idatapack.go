package ziface

/*
	封包、拆包 模块
	直接面向 TCP 连接中的数据流，用于处理TCP “粘包” 问题
*/

type IDataPack interface {
	// GetHeadLen 获取包头长度方法
	GetHeadLen() uint32
	// Pack 封包方法
	Pack(msg IMessage) ([]byte, error)
	// UnPack 拆包方法
	UnPack([]byte) (IMessage, error)
}
