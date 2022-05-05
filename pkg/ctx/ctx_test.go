package ctx_test

import (
	"context"
	"testing"
	"time"
)

func cleanup (cancel context.CancelFunc, resChan chan<- int) {
	cancel()
	close(resChan)
}

func TestCancellingChildDoesNotCancelParent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fooBar(ctx, 1)
	helloWorld(ctx, 20)

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
	ctx, cancel := context.WithCancel(context.Background())

	resFooBarC := fooBar(ctx, 10)
	resHelloWorldC := helloWorld(ctx, 20)

	time.Sleep(5 * time.Second)
	cancel()

	expectedMaxFooBar := 10
	if resFooBar := <-resFooBarC; expectedMaxFooBar == resFooBar {
		t.Fatalf("%d == %d", expectedMaxFooBar, resFooBar)
	}

	expectedMaxHelloWorld := 20
	if resHelloWorld := <-resHelloWorldC; expectedMaxHelloWorld == resHelloWorld {
		t.Fatalf("%d == %d", expectedMaxHelloWorld, resHelloWorld)
	}
}

func fooBar(ctx context.Context, tries int) <-chan int {
	resC := make(chan int)
	go func(resC chan int) {
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
	}(resC)

	return resC
}

func helloWorld(ctx context.Context, tries int) <-chan int {
	resC := make(chan int)
	go func(resC chan int) {
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
	}(resC)

	return resC
}
