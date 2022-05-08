package handler

import (
	"testing"

	"github.com/bedakb/codecrafters-redis/parser"
	"github.com/bedakb/codecrafters-redis/store"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := store.New()
	h := New(s)
	require.NotNil(t, h)
}

func TestHandle_ErrorIfGivenValueIsNotArray(t *testing.T) {
	h := &Handler{storage: store.New()}
	v := "$5\r\nget\r\n"
	got, err := h.Handle([]byte(v))
	require.Error(t, err)
	require.Equal(t, "", got)
}

func TestHandle_ErrorIfGivenValueIsEmptyArray(t *testing.T) {
	h := &Handler{storage: store.New()}
	v := "*0\r\n"
	got, err := h.Handle([]byte(v))
	require.Error(t, err)
	require.Equal(t, "", got)
}

func Test_handlePingCmd(t *testing.T) {
	w := "+PONG\r\n"
	got, err := handlePingCmd()
	require.NoError(t, err)
	require.Equal(t, w, got)
}

func Test_handleEchoCmd(t *testing.T) {
	w := "$5\r\nhello\r\n"
	args := []parser.Result{
		{
			Type:  parser.RedisBulkString,
			Value: "hello",
		},
	}
	got, err := handleEchoCmd(args)
	require.NoError(t, err)
	require.Equal(t, w, got)
}

func Test_handleEchoCmd_ErrorIfArgsAreMissing(t *testing.T) {
	args := []parser.Result{}
	got, err := handleEchoCmd(args)
	require.Error(t, err)
	require.Equal(t, "", got)
}

func Test_handleGetCmd(t *testing.T) {
	w := "$3\r\nbar\r\n"
	s := store.New()
	args := []parser.Result{
		{
			Type:  parser.RedisBulkString,
			Value: "foo",
		},
	}
	s.Set("foo", "bar")
	got, err := handleGetCmd(args, s)
	require.NoError(t, err)
	require.Equal(t, w, got)
}

func Test_handleGetCmd_ErrorIfArgsAreNotExactlyOne(t *testing.T) {
	s := store.New()
	args := []parser.Result{
		{
			Type:  parser.RedisBulkString,
			Value: "foo",
		},
		{
			Type:  parser.RedisBulkString,
			Value: "bar",
		},
	}
	s.Set("foo", "bar")
	got, err := handleGetCmd(args, s)
	require.Error(t, err)
	require.Equal(t, "", got)
}

func Test_handleSetCmd(t *testing.T) {
	w := "+OK\r\n"
	s := store.New()
	args := []parser.Result{
		{
			Type:  parser.RedisBulkString,
			Value: "foo",
		},
		{
			Type:  parser.RedisBulkString,
			Value: "bar",
		},
	}
	got, err := handleSetCmd(args, s)
	require.NoError(t, err)
	require.Equal(t, w, got)
}

func Test_handleSetCmd_ErrorIfArgAreNotSufficient(t *testing.T) {
	s := store.New()
	args := []parser.Result{
		{
			Type:  parser.RedisBulkString,
			Value: "foo",
		},
	}
	got, err := handleSetCmd(args, s)
	require.Error(t, err)
	require.Equal(t, "", got)
}

func Test_handleSetCmd_WithPxOption(t *testing.T) {
	w := "+OK\r\n"
	s := store.New()
	args := []parser.Result{
		{
			Type:  parser.RedisBulkString,
			Value: "foo",
		},
		{
			Type:  parser.RedisBulkString,
			Value: "bar",
		},
		{
			Type:  parser.RedisBulkString,
			Value: "px",
		},
		{
			Type:  parser.RedisBulkString,
			Value: "100",
		},
	}
	got, err := handleSetCmd(args, s)
	require.NoError(t, err)
	require.Equal(t, w, got)
}
