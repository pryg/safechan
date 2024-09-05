package safechan

import (
	"context"
)

type SafeChan[T any] struct {
	channel chan T
}

func New[T any]() SafeChan[T] {
	return SafeChan[T]{
		channel: make(chan T),
	}
}

func NewWithBuffer[T any](size uint) SafeChan[T] {
	return SafeChan[T]{
		channel: make(chan T, size),
	}
}

func (c SafeChan[T]) Send(ctx context.Context, m T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.channel <- m:
		return nil
	}
}

func (c SafeChan[T]) Receive(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		return *new(T), ctx.Err()
	case m := <-c.channel:
		return m, nil
	}
}

func (c SafeChan[T]) Merge(ctx context.Context, from SafeChan[T]) {
	go c.resend(ctx, from)
}

func (c SafeChan[T]) resend(ctx context.Context, from SafeChan[T]) {
	for {
		message, err := from.Receive(ctx)
		if err != nil {
			return
		}
		if err := c.Send(ctx, message); err != nil {
			return
		}
	}
}
