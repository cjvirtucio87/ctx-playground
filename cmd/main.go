package main

import (
	"context"
	"fmt"
	"time"
)

func cancellingParentCancelsChildren() {
	fmt.Println("cancelling parent cancels children")
	ctx, cancel := context.WithCancel(context.Background())
	defer func(cancel context.CancelFunc) {
		cancel()
		fmt.Println("done")
	}(cancel)

	go fooBar(ctx, 20)
	go helloWorld(ctx, 20)

	time.Sleep(5 * time.Second)
}

func fooBar(ctx context.Context, tries int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := tries; i >= 0; i-- {
		printMsg("fooBar", "foo!", tries - i)
		time.Sleep(1 * time.Second)
	}
}

func helloWorld(ctx context.Context, tries int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := tries; i >= 0; i-- {
		printMsg("helloWorld", "hello!", tries - i)
		time.Sleep(1 * time.Second)
	}
}

func printMsg(src string, msg string, num int) {
	fmt.Printf("[%s | %d] %s\n", src, num, msg)
}

func main() {
	cancellingParentCancelsChildren()
}
