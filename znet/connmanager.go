package znet

import (
	"go-ros-fog/ziface"
	"errors"
	"fmt"
	"sync"
)

// 连接管理模块

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 保护连接集合的读写锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将 conn 加入到 ConnManager 中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID= ", conn.GetConnID(), "add to ConnManager successfully: conn num = ", connMgr.Count())
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID= ", conn.GetConnID(), "remove to ConnManager successfully: conn num = ", connMgr.Count())
}

// Get 根据 connID 获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not Found")
}

// Count 得到当前连接总数
func (connMgr *ConnManager) Count() int {
	return len(connMgr.connections)
}

// ClearConn 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源 map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除并停止 conn 的工作
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections success! conn num = ", connMgr.Count())
}
