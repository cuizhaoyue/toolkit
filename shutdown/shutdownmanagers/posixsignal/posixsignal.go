package posixsignal

import (
	"github.com/cuizhaoyue/toolkit/shutdown"
	"os"
	"os/signal"
	"syscall"
)

// Name 定义shutdown manager的名称
const Name = "PosixSignalManager"

// PosixSignalManager 实现了ShutdownManager接口，它被添加到GracefullShutdown中
// 使用NewPosixSignalManager初始化
type PosixSignalManager struct {
	signals []os.Signal
}

var _ shutdown.ShutdownManager = (*PosixSignalManager)(nil)

// NewPosixSignalManager 初始化PosixSignalManager
// 可以提供用于监听的系统信号作为参数，如果没有提供则默认使用 SIGINT 和 SIGTERM
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signals: sig,
	}
}

// GetName 返回ShutdownManager的名称
func (p *PosixSignalManager) GetName() string {
	return Name
}

// Start 开始监听 posix 信号
func (p *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, p.signals...)

		// 接收到系统信号前阻塞
		<-c

		gs.StartShutdown(p)
	}()

	return nil
}

func (p *PosixSignalManager) ShutdownStart() error {
	return nil
}

func (p *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}
