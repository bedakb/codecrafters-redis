package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_decodeString(t *testing.T) {
	s := "+PING\r\n"
	w := Result{
		Type:  RedisSimpleString,
		Value: "PING",
	}

	got, gotPos := decodeString([]byte(s))
	require.Equal(t, w, got)
	require.Equal(t, 7, gotPos)
}

func Test_decodeError(t *testing.T) {
	s := "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
	w := Result{
		Type:  RedisError,
		Value: "WRONGTYPE Operation against a key holding the wrong kind of value",
	}

	got, gotPos := decodeError([]byte(s))
	require.Equal(t, w, got)
	require.Equal(t, 68, gotPos)
}

func Test_decodeInt(t *testing.T) {
	s := ":42\r\n"
	w := Result{
		Type:  RedisInt,
		Value: 42,
	}

	got, gotPos := decodeInt([]byte(s))
	require.Equal(t, w, got)
	require.Equal(t, 5, gotPos)
}

func Test_decodeInt_RedisErrorForNonInts(t *testing.T) {
	s := ":foo\r\n"
	w := Result{
		Type:  RedisError,
		Value: "cannot convert foo to int",
	}

	got, gotPos := decodeInt([]byte(s))
	require.Equal(t, w, got)
	require.Equal(t, 6, gotPos)
}

func Test_decodeBulkString(t *testing.T) {
	s := "$6\r\nhello\r\n"
	got, gotPos := decodeBulkString([]byte(s))
	require.Equal(t, "hello", got.Value)
	require.Equal(t, 11, gotPos)
}

func Test_decodeBulkString_RedisErrorWhenStringLenIsNotValid(t *testing.T) {
	s := "$baz\r\nhello\r\n"
	w := Result{
		Type:  RedisError,
		Value: "invalid bulk string size",
	}
	got, gotPos := decodeBulkString([]byte(s))
	require.Equal(t, w, got)
	require.Equal(t, 6, gotPos)
}

func Test_decodeArray(t *testing.T) {
	s := "*3\r\n$3\r\nhello\r\n$3\r\nworld\r\n:42\r\n"
	got, gotPos := decodeArray([]byte(s))
	want := Result{
		Type: RedisArray,
		Value: []Result{
			{
				Type:  RedisBulkString,
				Value: "hello",
			},
			{
				Type:  RedisBulkString,
				Value: "world",
			},
			{
				Type:  RedisInt,
				Value: 42,
			},
		},
	}
	require.Equal(t, want, got)
	require.Equal(t, 31, gotPos)
}

func Test_decodeArray_RedisErrorWhenValueIndicatingArrLenIsInvalid(t *testing.T) {
	s := "*foo\r\n$3\r\nhello\r\n$3\r\nworld\r\n:42\r\n"
	got, gotPos := decodeArray([]byte(s))
	want := Result{
		Type:  RedisError,
		Value: "invalid array size",
	}
	require.Equal(t, want, got)
	require.Equal(t, 6, gotPos)
}

func Test_matchCRLF(t *testing.T) {
	s := "wo\r\nrld"
	i, end := matchCRLF([]byte(s))
	require.Equal(t, 2, i)
	require.Equal(t, 4, end)
}
