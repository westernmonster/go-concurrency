package split

import "context"

func Split(ctx context.Context, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)

		for v := range in {
			select {
			case <-ctx.Done():
				return
			case out1 <- v:
			}

			select {
			case <-ctx.Done():
				return
			case out2 <- v:
			}
		}
	}()

	return out1, out2
}
