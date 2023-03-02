package shutdown

import "sync"

// ShutdownCallback 是一个为回调必须实现的接口
// 当请求关机时，Onshutdown 函数将会被调用，参数是请求关机的ShutdownManager的名称
type ShutdownCallback interface {
	OnShutdown(string) error
}

// ShutdownFunc 是一个帮助类型，用于提供作为ShutdownCallback的匿名函数
type ShutdownFunc func(string) error

// OnShutdown 定义当关机操作被触发时需要去运行的动作
func (f ShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

// ShutdownManager 是一个由ShutdownManager实现的接口
type ShutdownManager interface {
	// GetName 获取ShutdownManager的名称
	GetName() string
	// Start ShutdownManagers开始监听关机请求
	Start(gs GSInterface) error
	// ShutdownStart 当GSInterface中的StartShutdown被调用时，ShutdownStart()会被调用，
	// 然后所有的ShutdownCallback都被执行，一旦所有的ShutdownCallback返回，ShutdownFinish会被调用
	ShutdownStart() error
	ShutdownFinish() error
}

// ErrorHandler 是一个可以传给SetErrorHander来处理异步错误的接口
type ErrorHandler interface {
	OnError(error)
}

// ErrorFunc 是一个helper类型，你可以提供作为ErrorHandler的匿名函数
type ErrorFunc func(error)

// OnError 定义了当错误发生时需要执行的动作
func (f ErrorFunc) OnError(err error) {
	f(err)
}

// GSInterface 是一个由GracefulShutdown实现的接口，
// 当请求关机时它会传递给ShutdownManager然后调用StartShudown
type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback ShutdownCallback)
}

// GracefulShutdown 是处理ShutdownCallback和ShutdownManager的主要结构体。使用New初始化它
type GracefulShutdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}

// New 初始化一个GracefulShutdown
func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 0, 10),
		managers:  make([]ShutdownManager, 0, 3),
	}
}

// Start 调用所有添加的ShutdownManager的Start方法，ShutdownManager开始监听关机请求
// 如果ShutdowManager返回错误则返回该错误
func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}

	return nil
}

// AddShutdownManager 添加一个即将监听关机请求的ShutdownManager实例
func (gs *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}

// AddShutdownCallback 添加一个在关机的时候被调用的ShutdownCallback实例
//
// 也可以添加任何一个实现了ShutdownCallback接口的对象，或者可以使用一个类似这种的函数：
//
//	AddShutdownCallback(shutdown.ShutdownFunc(func() error {
//	     // callback code
//		    return nil
//	 }))
func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

// SetErrorHandler 设置一个当在ShutdownCallback或ShutdownManager中发生错误时要调用的ErrorHandler实例
//
// 也可以提供任何实现了ErrorHandler接口的对象，或者提供一种类似这种的函数：
//
//	SetErrorHandler(shutdown.ErrorFunc(func (err error) {
//		// handle error
//	}))
func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

// ReportError 是一个给ErrorHandler汇报错误的函数，它在ShutdownManager中使用
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}

// StartShutdown 从ShutdownManager中被调用，并且开始关机
// 首先调用ShutdownManager中的ShutdownStart
// 然后调用所有的ShutdownCallback，等待所有的callback完成
// 最后调用ShutdownManager中的ShutdownFinish
func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(callback ShutdownCallback) {
			defer wg.Done()

			gs.ReportError(callback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}

	wg.Wait()

	gs.ReportError(sm.ShutdownFinish())
}
