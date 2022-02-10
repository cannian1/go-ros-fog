package znet

import "go-ros-fog/ziface"

// BaseRouter 实现router时，嵌入这个基类，根据需要对这个基类的方法重写
type BaseRouter struct {
}

// 这里 BaseRouter 的方法都为空
// 是因为有的 Router 不希望有 PreHandle、PostHandle 这两个业务
// 所以 Router 全部继承 BaseRouter的好处是不需要实现 PreHandle、PostHandle

// PreHandle 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

// Handle 处理conn业务的主方法hook
func (br *BaseRouter) Handle(request ziface.IRequest) {

}

// PostHandle 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
