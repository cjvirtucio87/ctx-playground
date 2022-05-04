package main

import (
	"context"
	"fmt"
	"time"
)

func cancelWithMsg(cancel context.CancelFunc, src string) {
	defer cancel()
	fmt.Printf("[%s] cancelling\n", src)
}

func cancellingChildDoesNotCancelParent() {
	fmt.Println("cancelling child does not cancel parent")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancelWithMsg(cancel, "cancellingChildDoesNotCancelParent")

	go fooBar(ctx, 2)
	go helloWorld(ctx, 20)

	time.Sleep(10 * time.Second)
}

func cancellingParentCancelsChildren() {
	fmt.Println("cancelling parent cancels children")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancelWithMsg(cancel, "cancellingParentCancelsChildren")

	go fooBar(ctx, 20)
	go helloWorld(ctx, 20)

	time.Sleep(5 * time.Second)
}

func fooBar(ctx context.Context, tries int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancelWithMsg(cancel, "fooBar")

	for i := tries; i >= 0; i-- {
		select {
		case <- ctx.Done():
			return
		default:
			printMsg("fooBar", "foo!", tries - i)
			time.Sleep(1 * time.Second)
		}
	}
}

func helloWorld(ctx context.Context, tries int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancelWithMsg(cancel, "helloWorld")

	for i := tries; i >= 0; i-- {
		select {
		case <- ctx.Done():
			return
		default:
			printMsg("helloWorld", "hello!", tries - i)
			time.Sleep(1 * time.Second)
		}
	}
}

func printMsg(src string, msg string, num int) {
	fmt.Printf("[%s | %d] %s\n", src, num, msg)
}

func main() {
	cancellingParentCancelsChildren()
	time.Sleep(5 * time.Second)

	cancellingChildDoesNotCancelParent()
	time.Sleep(5 * time.Second)
}
