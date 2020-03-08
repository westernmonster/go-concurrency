package take

import "context"

func Repeat(ctx context.Context, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for _, v := range values {
			select {
			case <-ctx.Done():
				return
			case valueStream <- v:
			}
		}
	}()
	return valueStream
}

func RepeatFn(ctx context.Context, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-ctx.Done():
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

func Take(ctx context.Context, in <-chan interface{}, n int) (out <-chan interface{}) {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case valueStream <- <-in:
			}
		}
	}()

	return valueStream
}
