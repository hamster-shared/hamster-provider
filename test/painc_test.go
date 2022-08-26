package test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
)

func TestPainc(t *testing.T) {
	main()
}

func main() {

	fmt.Println("start main")
	ctx, _ := context.WithCancel(context.Background())
	isPanic := make(chan bool)

	go start_listener(ctx, isPanic)

	select {}
}

func start_listener(ctx context.Context, isPanic chan bool) {
	for {
		go watchEvent(ctx, isPanic)
		select {
		case <-isPanic:
			go watchEvent(ctx, isPanic)
		}
	}
}

func watchEvent(ctx context.Context, channel chan bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			channel <- true
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		default:
			i := decodeSomeEvent()
			fmt.Println("deal event :", i)
		}

	}

}

func decodeSomeEvent() int {

	i := rand.Intn(10)
	if i > 8 {
		panic("some eror")
	}
	return i
}
