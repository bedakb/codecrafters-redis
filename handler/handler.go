package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bedakb/codecrafters-redis/parser"
	"github.com/bedakb/codecrafters-redis/store"
)

const (
	PingCmd = "ping"
	EchoCmd = "echo"
	GetCmd  = "get"
	SetCmd  = "set"
)

// PxOption is an option used to set key expiration.
const PxOption = "px"

// Handler handles input sent by redis client.
type Handler struct {
	storage *store.Store
}

// New creates a new instance of the Handler.
func New(storage *store.Store) *Handler {
	return &Handler{
		storage: storage,
	}
}

// Handle reads the value sent by redis client and executes necessary command.
//
// It will return an encoded response ready to be read by the client and the error.
func (h *Handler) Handle(value []byte) (string, error) {
	v := parser.Decode(value)
	if v.Type != parser.RedisArray {
		return "", errors.New("redis command must be sent as an array")
	}

	res, ok := v.Value.([]parser.Result)
	if !ok {
		return "", errors.New("cannot convert value into list of the results")
	}
	if len(res) == 0 {
		return "", errors.New("command and arguments are empty")
	}

	cmd := res[0]
	args := res[1:]

	return h.handleCommand(cmd, args)
}

func (h *Handler) handleCommand(cmd parser.Result, args []parser.Result) (string, error) {
	switch cmd.Value {
	case PingCmd:
		return handlePingCmd()
	case EchoCmd:
		return handleEchoCmd(args)
	case GetCmd:
		return handleGetCmd(args, h.storage)
	case SetCmd:
		return handleSetCmd(args, h.storage)
	}
	return "", fmt.Errorf("unknown command received")
}

func handlePingCmd() (string, error) {
	return parser.Encode(parser.Result{
		Type:  parser.RedisSimpleString,
		Value: "PONG",
	})
}

func handleEchoCmd(args []parser.Result) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("wrong number of arguments: got 0, expected at least 1")
	}

	vals := make([]interface{}, 0, len(args))
	for _, a := range args {
		vals = append(vals, a.Value)
	}

	v, err := stringify(vals)
	if err != nil {
		return "", fmt.Errorf("cannot read the args: %w", err)
	}

	r := parser.Result{
		Type:  parser.RedisBulkString,
		Value: strings.Join(v, " "),
	}

	return parser.Encode(r)
}

func handleGetCmd(args []parser.Result, s *store.Store) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("wrong number of arguments: got %d, expected 1", len(args))
	}

	k, ok := args[0].Value.(string)
	if !ok {
		return "", errors.New("key is not a string")
	}

	r := parser.Result{
		Type: parser.RedisBulkString,
	}

	v, ok := s.Get(k)
	if ok {
		r.Value = v
	}

	return parser.Encode(r)
}

func handleSetCmd(args []parser.Result, s *store.Store) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("wrong number of arguments: got %d, expected at least 2", len(args))
	}

	k := args[0].Value.(string)
	v := args[1].Value.(string)

	opts := args[2:]
	if setCmdContainsOpts(opts) {
		o, dur, err := extractSetOptions(opts)
		if err != nil {
			return "", fmt.Errorf("cannot extract options from the set command: %w", err)
		}

		switch o {
		case PxOption:
			s.SetWithExpirity(k, v, time.Millisecond*dur)
		default:
		}
	} else {
		s.Set(k, v)
	}

	r := parser.Result{
		Type:  parser.RedisSimpleString,
		Value: "OK",
	}
	return parser.Encode(r)
}

func extractSetOptions(options []parser.Result) (string, time.Duration, error) {
	o, ok := options[0].Value.(string)
	if !ok {
		return "", 0, errors.New("option is not a string")
	}

	d, ok := options[1].Value.(string)
	if !ok {
		return "", 0, errors.New("reading duration option failed")
	}

	dur, err := strconv.Atoi(d)
	if err != nil {
		return "", 0, errors.New("duration must be a number")
	}

	return o, time.Duration(dur), nil
}

func setCmdContainsOpts(args []parser.Result) bool {
	return len(args) >= 1
}

func stringify(values []interface{}) ([]string, error) {
	s := make([]string, 0, len(values))
	for _, v := range values {
		sv, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("cannot convert %v to string", sv)
		}
		s = append(s, sv)
	}
	return s, nil
}
