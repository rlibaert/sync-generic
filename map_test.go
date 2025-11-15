package sync_test

import (
	"testing"

	"github.com/rlibaert/sync-generic"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	t.Run("Store & Load", func(t *testing.T) {
		var m sync.Map[string, int]
		m.Store("foo", 42)
		value, ok := m.Load("foo")
		require.True(t, ok)
		require.Equal(t, 42, value)
	})

	t.Run("Delete", func(t *testing.T) {
		var m sync.Map[string, int]
		m.Store("foo", 42)
		m.Delete("foo")
		value, ok := m.Load("foo")
		require.False(t, ok)
		require.Equal(t, 0, value)
	})

	t.Run("LoadOrStore", func(t *testing.T) {
		var m sync.Map[string, int]

		t.Run("stores", func(t *testing.T) {
			actual, loaded := m.LoadOrStore("foo", 42)
			require.False(t, loaded)
			require.Equal(t, 42, actual)
		})

		t.Run("loads", func(t *testing.T) {
			actual, loaded := m.LoadOrStore("foo", 12345)
			require.True(t, loaded)
			require.Equal(t, 42, actual)
		})
	})

	t.Run("LoadAndDelete", func(t *testing.T) {
		var m sync.Map[string, int]
		m.Store("foo", 42)

		t.Run("loads", func(t *testing.T) {
			value, loaded := m.LoadAndDelete("foo")
			require.True(t, loaded)
			require.Equal(t, 42, value)
		})

		t.Run("is deleted", func(t *testing.T) {
			value, loaded := m.LoadAndDelete("foo")
			require.False(t, loaded)
			require.Equal(t, 0, value)
		})
	})

	t.Run("Range", func(t *testing.T) {
		var m sync.Map[string, int]
		m.Store("foo", 42)
		m.Store("bar", 12)
		m.Store("baz", 69)

		ranged := map[string]int{}
		m.Range(func(key string, value int) bool {
			ranged[key] = value
			return true
		})

		require.Len(t, ranged, 3)
		require.Equal(t, 42, ranged["foo"])
		require.Equal(t, 12, ranged["bar"])
		require.Equal(t, 69, ranged["baz"])
	})

	t.Run("Clear", func(t *testing.T) {
		var m sync.Map[string, int]
		m.Store("foo", 42)
		m.Store("bar", 12)
		m.Store("baz", 69)
		m.Clear()

		ranged := map[string]int{}
		m.Range(func(key string, value int) bool {
			ranged[key] = value
			return true
		})

		require.Empty(t, ranged)
	})

	t.Run("Swap", func(t *testing.T) {
		t.Run("loads", func(t *testing.T) {
			var m sync.Map[string, int]
			m.Store("foo", 42)

			previous, loaded := m.Swap("foo", 12)
			require.True(t, loaded)
			require.Equal(t, 42, previous)

			value, ok := m.Load("foo")
			require.True(t, ok)
			require.Equal(t, 12, value)
		})

		t.Run("does not load", func(t *testing.T) {
			var m sync.Map[string, int]

			previous, loaded := m.Swap("foo", 12)
			require.False(t, loaded)
			require.Equal(t, 0, previous)

			value, ok := m.Load("foo")
			require.True(t, ok)
			require.Equal(t, 12, value)
		})
	})

	t.Run("CompareAndSwap", func(t *testing.T) {
		var m sync.Map[string, int]
		m.Store("foo", 42)

		t.Run("does not swap", func(t *testing.T) {
			require.False(t, m.CompareAndSwap("foo", 12, 69))
			value, ok := m.Load("foo")
			require.True(t, ok)
			require.Equal(t, 42, value)
		})

		t.Run("swaps", func(t *testing.T) {
			require.True(t, m.CompareAndSwap("foo", 42, 12))
			value, ok := m.Load("foo")
			require.True(t, ok)
			require.Equal(t, 12, value)
		})
	})

	t.Run("CompareAndDelete", func(t *testing.T) {
		var m sync.Map[string, int]
		m.Store("foo", 42)

		t.Run("does not delete", func(t *testing.T) {
			require.False(t, m.CompareAndDelete("foo", 12))
			value, ok := m.Load("foo")
			require.True(t, ok)
			require.Equal(t, 42, value)
		})

		t.Run("deletes", func(t *testing.T) {
			require.True(t, m.CompareAndDelete("foo", 42))
			value, ok := m.Load("foo")
			require.False(t, ok)
			require.Equal(t, 0, value)
		})
	})
}

func TestMapPanics(t *testing.T) {
	t.Run("CompareAndSwap", func(t *testing.T) {
		require.PanicsWithError(t, "runtime error: comparing uncomparable type []uint8", func() {
			var m sync.Map[string, []byte]
			m.Store("foo", []byte("hello"))
			m.CompareAndSwap("foo", []byte("hello"), []byte("world"))
		})
	})

	t.Run("CompareAndDelete", func(t *testing.T) {
		require.PanicsWithError(t, "runtime error: comparing uncomparable type []uint8", func() {
			var m sync.Map[string, []byte]
			m.Store("foo", []byte("hello"))
			m.CompareAndDelete("foo", []byte("hello"))
		})
	})
}
