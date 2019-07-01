package RContext

// 限制最大执行数量
// Error后不执行后续任务，
// 无Error时Wait等待运行完成

import (
	"context"
	"fmt"
	"sync"
)

type RContext struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup // wait
	n      chan int        // 控制个数
}

func NewRContext(ctx context.Context, threadNumMax int) *RContext {
	t := new(RContext)
	t.ctx, t.cancel = context.WithCancel(ctx)
	t.n = make(chan int, threadNumMax)
	t.wg = &sync.WaitGroup{}
	return t
}

func (t *RContext) Add() (e error) {
	select {
	case <- t.ctx.Done():
		return fmt.Errorf("context canceled")
	default:
	}

	t.n <- 1
	t.wg.Add(1)
	return
}
func (t *RContext) Done() {
	<-t.n
	t.wg.Done()
}
func (t *RContext) Wait() {
	t.wg.Wait()
}

func (t *RContext) Context() context.Context {
	return t.ctx
}
func (t *RContext) Cancel() {
	t.cancel()
}