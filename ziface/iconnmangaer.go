package ziface

// 连接管理模块抽象层

type IConnManager interface {
	// Add 添加连接
	Add(conn IConnection)
	// Remove 删除连接
	Remove(conn IConnection)
	// Get 根据 connID 获取连接
	Get(connID uint32) (IConnection, error)
	// Count 得到当前连接总数
	Count() int
	// ClearConn 清除并终止所有连接
	ClearConn()
}
