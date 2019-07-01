package RContext

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNewRoutineContext1(t *testing.T) {
	rctxl := NewRContext(context.Background(), 2)

	for i := 0; i < 7; i++ {
		if rctxl.Add() != nil {
			break
		}
		go func() {
			routineContext1Func1(rctxl, i)
			rctxl.Cancel()
		}()
	}
	rctxl.Wait()
	fmt.Println("rctxl.Wait()")

	time.Sleep(5 * time.Second)
}
func routineContext1Func1(rctx *RContext, i int) (e error) {
	fmt.Println(i, "routineContext1Func1 in")
	defer rctx.Done()

	select {
	case <- rctx.Context().Done():
		return
	default:
	}

	time.Sleep(1 * time.Second)
	rctx2 := NewRContext(rctx.Context(), 2)
	for i := 0; i < 7; i++ {
		if e := rctx2.Add(); e != nil {
			break
		}
		go func() {
			e := routineContext1Func2(rctx2, i)
			if e != nil {
				rctx2.cancel()
			}
		}()
	}
	rctx2.Wait()
	fmt.Println("rctx2.Wait()")
	time.Sleep(1 * time.Second)

	fmt.Println(i, "routineContext1Func1 out")
	return
}
func routineContext1Func2(rctx *RContext, i int) (e error) {
	fmt.Println(i, "routineContext1Func2 in")
	defer rctx.Done()

	select {
	case <- rctx.Context().Done():
		return
	default:
	}

	time.Sleep(1 * time.Second)
	if i == 2 {
		e = errors.New("routineContext1Func2 out send error")
		fmt.Println(2, "routineContext1Func2 out send error")
		return
	}
	time.Sleep(1 * time.Second)

	fmt.Println(i, "routineContext1Func2 out")
	return
}