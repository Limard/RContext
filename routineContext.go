package RContext

// 限制最大执行数量
// Error后不执行后续任务，
// 无Error时Wait等待运行完成

import (
	"context"
	"errors"
	"time"
)

type RoutineContext struct {
	ctx    context.Context
	cancel context.CancelFunc
	n      chan int
	err    error
}

func NewRoutineContext(ctx context.Context, threadNumMax int) *RoutineContext {
	t := new(RoutineContext)
	t.ctx, t.cancel = context.WithCancel(ctx)
	t.n = make(chan int, threadNumMax)
	return t
}

func (t *RoutineContext) Add() error {
	if t.err != nil {
		return t.err
	}
	t.n <- 1
	return nil
}
func (t *RoutineContext) CheckError() error {
	return t.err
}
func (t *RoutineContext) Done() {
	<-t.n
	if len(t.n) == 0 {
		t.cancel()
	}
}
func (t *RoutineContext) Wait() error {
	select {
	case <-t.ctx.Done():
		return t.err
	}
	return nil
}
func (t *RoutineContext) WaitTimeout(timeout time.Duration) error {
	select {
	case <-t.ctx.Done():
		return t.err
	case <-time.After(timeout):
		return errors.New("timeout")
	}
	return nil
}
func (t *RoutineContext) Error(err error) {
	t.err = err
	t.cancel()
}
func (t *RoutineContext) Context() context.Context {
	return t.ctx
}