package utils

import "time"

func DebouncerChannel[T any](in <-chan T, delay time.Duration, bufferSize int) <-chan []T {
	out := make(chan []T, bufferSize)

	go func() {
		buffers := make([]T, 0, bufferSize)
		isClosed := false

		for !isClosed {
			wait := false

		untilEmpty:
			for {
				select {
				case elem, ok := <-in:
					if !ok {
						for elem := range in {
							buffers = append(buffers, elem)
						}
						isClosed = true
						break untilEmpty
					} else {
						if !wait {
							time.Sleep(delay)
							wait = true
						}

						buffers = append(buffers, elem)
					}
				default:
					break untilEmpty
				}
			}

			if len(buffers) > 0 {
				out <- buffers
				buffers = make([]T, 0, bufferSize)
			}
		}

		close(out)
	}()

	return out
}
