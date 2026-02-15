package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const kvSeparator = ":"
const infoSeparator = "|"
const StateByteLength = 32

var ErrInvalidStateInfo = errors.New("invalid state info")

type State struct {
	kv    map[string]string
	state []byte
}

func (s *State) AddInfo(key, value string) {
	if key == "" || value == "" {
		return
	}

	s.kv[key] = url.QueryEscape(value)
}

func (s *State) GetInfo(key string) (string, bool) {
	if key == "" {
		return "", false
	}

	value, ok := s.kv[key]
	if !ok {
		return "", false
	}

	unescaped, err := url.QueryUnescape(value)
	if err != nil {
		return "", false
	}

	return unescaped, true
}

func (s *State) String() string {
	info := make([]string, 0, len(s.kv))
	for k, v := range s.kv {
		info = append(info, fmt.Sprintf("%s%s%s", k, kvSeparator, v))
	}

	infoStr := fmt.Sprintf("%s%s", s.state, strings.Join(info, infoSeparator))

	return base64.URLEncoding.EncodeToString([]byte(infoStr))
}

func NewState() (*State, error) {
	b := make([]byte, StateByteLength)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return &State{
		state: b,
		kv:    make(map[string]string),
	}, nil
}

func NewStateFromEncode(stateStr string) (*State, error) {
	b, err := base64.URLEncoding.DecodeString(stateStr)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(b), infoSeparator)

	if len(parts) < 1 {
		return nil, ErrInvalidStateInfo
	}

	state := State{
		state: []byte(parts[0]),
		kv:    make(map[string]string),
	}

	for _, k := range parts[1:] {
		kv := strings.Split(k, kvSeparator)

		if len(kv) != 2 { //nolint:mnd // checking for key and value parts
			return nil, ErrInvalidStateInfo
		}

		key := kv[0]
		value := kv[1]
		state.kv[key] = value
	}

	return &state, nil
}
