package safechan

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSafeChan_Send(t *testing.T) {
	t.Run("success, send to chan without buffer", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

		var (
			wg  sync.WaitGroup
			got []string
		)

		c := New[string]()
		receive := func() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				m, err := c.Receive(ctx)
				assert.NoError(t, err)
				got = append(got, m)
			}()
		}

		receive()
		assert.NoError(t, c.Send(ctx, "msg1"))
		wg.Wait()
		assert.Equal(t, []string{"msg1"}, got)

		receive()
		assert.NoError(t, c.Send(ctx, "msg2"))
		wg.Wait()
		assert.Equal(t, []string{"msg1", "msg2"}, got)

		cancel()
		err := c.Send(ctx, "foobar")
		assert.ErrorIs(t, context.Canceled, err)
		assert.Equal(t, []string{"msg1", "msg2"}, got)
	})

	t.Run("failed, send to chan without buffer error: context deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()

		c := New[string]()
		err := c.Send(ctx, "msg")
		assert.ErrorIs(t, context.DeadlineExceeded, err)
	})

	t.Run("success, send to chan with buffer", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		var (
			wg  sync.WaitGroup
			got []string
		)

		c := NewWithBuffer[string](3)
		receive := func() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				m, err := c.Receive(ctx)
				assert.NoError(t, err)
				got = append(got, m)
			}()
		}

		assert.NoError(t, c.Send(ctx, "msg1"))
		assert.NoError(t, c.Send(ctx, "msg2"))
		assert.NoError(t, c.Send(ctx, "msg3"))
		receive()
		receive()
		receive()
		wg.Wait()

		cancel()
		err := c.Send(ctx, "foobar")
		assert.ErrorIs(t, context.Canceled, err)
		assert.Equal(t, []string{"msg1", "msg2", "msg3"}, got)
	})

	t.Run("failed, send to chan with buffer error: context deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()

		c := NewWithBuffer[string](2)
		assert.NoError(t, c.Send(ctx, "msg1"))
		assert.NoError(t, c.Send(ctx, "msg2"))
		err := c.Send(ctx, "msg3")
		assert.ErrorIs(t, context.DeadlineExceeded, err)
	})
}
