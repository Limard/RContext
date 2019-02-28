package RContext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewRoutineContext1(t *testing.T) {
	rl := NewRoutineContext(context.Background(), 2)

	for i := 0; i < 7; i++ {
		if rl.Add() != nil {
			break
		}
		go routineContext1Func1(rl, i)
	}
	fmt.Println("rl.Wait():", rl.Wait())

	time.Sleep(5 * time.Second)
}
func routineContext1Func1(rctx *RoutineContext, i int) {
	fmt.Println(i, "routineContext1Func1 in")
	defer rctx.Done()

	time.Sleep(1 * time.Second)
	rctx2 := NewRoutineContext(context.Background(), 2)
	for i := 0; i < 7; i++ {
		if e := rctx2.Add(); e != nil {
			rctx.Error(e)
			break
		}
		go routineContext1Func2(rctx2, i)
	}
	fmt.Println("rctx2.Wait():", rctx2.Wait())
	time.Sleep(1 * time.Second)

	fmt.Println(i, "routineContext1Func1 out")
}
func routineContext1Func2(rctx *RoutineContext, i int) {
	fmt.Println(i, "routineContext1Func2 in")
	defer rctx.Done()

	time.Sleep(1 * time.Second)
	if i == 2 {
		rctx.Error(fmt.Errorf("error message"))
		fmt.Println(2, "routineContext1Func2 out send error")
		return
	}
	time.Sleep(1 * time.Second)

	fmt.Println(i, "routineContext1Func2 out")
}

func TestNewRoutineContext2(t *testing.T) {
	rctx := NewRoutineContext(context.Background(), 5)

	for i := 0; i < 7; i++ {
		if e := rctx.Add(); e != nil {
			fmt.Println("Main break", i)
			break
		}
		routineContext2Func1(rctx, i)
	}
	fmt.Println("rctx.Wait():", rctx.Wait())

	time.Sleep(5 * time.Second)
}
func routineContext2Func1(rctx *RoutineContext, i int) {
	fmt.Println("Func1 in", i)
	defer rctx.Done()

	time.Sleep(1 * time.Second)
	for i2 := 0; i2 < 7; i2++ {
		if e := rctx.Add(); e != nil {
			fmt.Println("Func1 break", i2)
			break
		}
		go routineContext2Func2(rctx, i2)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Func1 out", i)
}
func routineContext2Func2(rctx *RoutineContext, i int) {
	fmt.Println("Func2 in", i)
	defer rctx.Done()

	time.Sleep(1 * time.Second)
	if i == 2 {
		rctx.Error(fmt.Errorf("error message"))
		fmt.Println("Func2 out send error")
		return
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Func2 out", i)
}
