package healthcheck

import (
	"context"
	"gitlab.corp.qunar.com/tcdev/qconfig-go/common/logger"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"runtime"
	"time"
)

type Checker struct {
	state        *atomic.Bool  // 状态记录, bool类型, true表示健康状态, false表示非健康状态
	interval     time.Duration // 健康检测的间隔
	healthFunc   func() error  // 状态改变为 health 时执行此function
	unHealthFunc func() error  // 状态变为 unHealth 时执行此function

	ctx    context.Context
	cancel context.CancelFunc
	logger *zap.Logger
}

func NewChecker(interval time.Duration, healthFunc func() error, unHealthFunc func() error, logger *zap.Logger) *Checker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Checker{
		state:        atomic.NewBool(false),
		interval:     interval,
		healthFunc:   healthFunc,
		unHealthFunc: unHealthFunc,

		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
}

func (c *Checker) stateChanged(desired bool) bool {
	// 已经是期望的状态，直接返回false
	if c.state.Load() == desired {
		return false
	}

	for {
		// 如果状态与期望的不一直，则改变状态, 并返回true，说明状态有变化
		if c.state.CAS(c.state.Load(), desired) {
			return true
		}
		runtime.Gosched()
	}
}

// CheckForever 检测healthcheck, 允许指定 health 和 unHealth 时的hook 函数
func (c *Checker) CheckForever() {
	ticker := time.NewTicker(c.interval)

	for {
		select {
		case <-c.ctx.Done():
			logger.Info("[healthcheck] checker exit")
			return
		case <-ticker.C:
		}

		var err error
		if IsHealth() {
			if c.stateChanged(true) {
				if err = c.healthFunc(); err != nil {
					logger.Error("[healthcheck] run health function error ", zap.Error(err))
				}
			}
			continue
		}

		if c.stateChanged(false) {
			if err = c.unHealthFunc(); err != nil {
				logger.Error("[healthcheck] run unHealth function error ", zap.Error(err))
			}
		}
	}
}

// Cancel 关闭checker
func (c *Checker) Cancel() {
	c.cancel()
}
