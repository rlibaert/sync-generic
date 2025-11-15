package sync_test

import (
	"testing"

	"github.com/rlibaert/sync-generic"
	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	t.Run("returns zero value", func(t *testing.T) {
		p := new(sync.Pool[[]byte])
		require.Equal(t, []byte(nil), p.Get())
		p = p.New(nil)
		require.Equal(t, []byte(nil), p.Get())
	})

	t.Run("returns new value", func(t *testing.T) {
		p := new(sync.Pool[[]byte]).New(func() []byte { return make([]byte, 8) })
		x := p.Get()
		require.NotEqual(t, []byte(nil), x)
		require.Len(t, x, 8)
	})
}
