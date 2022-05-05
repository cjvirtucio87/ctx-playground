package ctx_test

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func cleanup (cancel context.CancelFunc, resChan chan<- int) {
	cancel()
	close(resChan)
}

func TestCancellingChildDoesNotCancelParent(t *testing.T) {
	fmt.Println("cancelling child does not cancel parent")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resFooBarC := make(chan int)
	go fooBar(ctx, 1, resFooBarC)

	resHelloWorldC := make(chan int)
	go helloWorld(ctx, 20, resHelloWorldC)

	expectedParentRes := 5
	var parentRes int
	for parentRes < expectedParentRes {
		parentRes += 1
		time.Sleep(1 * time.Second)
	}

	if expectedParentRes != parentRes {
		t.Fatalf("%d != %d", expectedParentRes, parentRes)
	}
}

func TestCancellingParentCancelsChildren(t *testing.T) {
	fmt.Println("cancelling parent cancels children")
	ctx, cancel := context.WithCancel(context.Background())
	resFooBarC := make(chan int)

	go fooBar(ctx, 10, resFooBarC)

	resHelloWorldC := make(chan int)
	go helloWorld(ctx, 20, resHelloWorldC)


	time.Sleep(5 * time.Second)
	cancel()

	expectedFooBar := 10
	if resFooBar := <-resFooBarC; expectedFooBar == resFooBar {
		t.Fatalf("%d == %d", expectedFooBar, resFooBar)
	}

	expectedHelloWorld := 20
	if resHelloWorld := <-resHelloWorldC; expectedHelloWorld == resHelloWorld {
		t.Fatalf("%d == %d", expectedHelloWorld, resHelloWorld)
	}
}

func fooBar(ctx context.Context, tries int, resC chan<- int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cleanup(cancel, resC)

	var val int
	for i := tries; i >= 0; i-- {
		val = tries - i
		select {
		case <- ctx.Done():
			resC <-val
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}

	resC <-val
}

func helloWorld(ctx context.Context, tries int, resC chan<- int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cleanup(cancel, resC)

	var val int
	for i := tries; i >= 0; i-- {
		val = tries - i
		select {
		case <- ctx.Done():
			resC <-val
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}

	resC <-val
}
