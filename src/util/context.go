package util

import (
	"context"
	"fmt"
)

type SimpleContext struct {
	ctx      context.Context
	cancelFn func()
}

func NewSimpleContext() *SimpleContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &SimpleContext{
		ctx:      ctx,
		cancelFn: cancel,
	}
}

func (sc *SimpleContext) Context() context.Context {
	return sc.ctx
}

func (sc *SimpleContext) Cancel() {
	sc.cancelFn()
}

func (sc *SimpleContext) IsCancelled() bool {
	select {
	case <-sc.ctx.Done():
		fmt.Println("취소 확인됩")
		return true
	default:
		return false
	}
}
