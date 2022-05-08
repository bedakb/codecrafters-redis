package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := New()
	require.NotNil(t, s)
}

func TestSet(t *testing.T) {
	s := New()

	s.Set("foo", "bar")
	s.Set("one", "two")

	require.Equal(t, 2, s.Len())
}

func TestSetWithExpirity(t *testing.T) {
	s := New()

	s.Set("foo", "bar")
	s.SetWithExpirity("one", "two", time.Millisecond*100)

	v, ok := s.Get("foo")
	require.Equal(t, "bar", v)
	require.Equal(t, true, ok)

	v, ok = s.Get("one")
	require.Equal(t, "two", v)
	require.Equal(t, true, ok)

	time.Sleep(time.Millisecond * 105)

	v, ok = s.Get("one")
	require.Equal(t, "", v)
	require.Equal(t, false, ok)
}

func TestGet_ExistingKey(t *testing.T) {
	s := New()

	s.Set("foo", "bar")

	v, ok := s.Get("foo")
	require.Equal(t, "bar", v)
	require.Equal(t, true, ok)
}

func TestGet_NonExistingKey(t *testing.T) {
	s := New()

	s.Set("foo", "bar")
	s.Set("one", "two")

	v, ok := s.Get("key")
	require.Equal(t, "", v)
	require.Equal(t, false, ok)
}
