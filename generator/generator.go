package generator

import "context"

func Generator(ctx context.Context, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-ctx.Done():
				return
			case intStream <- i:
			}
		}
	}()

	return intStream
}
