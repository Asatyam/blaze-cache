package store

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type MapToMap map[string]map[string]any

type Store struct {
	mu sync.Mutex
	mp MapToMap
}

func NewStore() *Store {
	return &Store{
		mp: MapToMap{},
	}
}

// Set arr would be of form ["sym", "key", "sym", "value", opts...  ]
// opts ["sym", "opt", "sym", "optValue",..........]

func (s *Store) Set(arrString []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := arrString[1]
	keyMap, ok := s.mp[key]
	oldValue := ""
	if !ok {
		s.mp[key] = make(map[string]any)
	} else {
		oldValue = keyMap["Value"].(string)
	}
	toReturnOld := false

	s.mp[key]["Value"] = arrString[3]
	for i := 5; i < len(arrString); {
		opt := strings.ToUpper(arrString[i])
		switch opt {
		case "GET":
			toReturnOld = true
			i += 2

		case "PX":
			expiryTime, err := s.GetExpiryTime(arrString[i+2], 1000000)
			if err != nil {
				return "", err
			}
			s.mp[key][opt] = expiryTime
			i += 4
		default:
			s.mp[key][opt] = arrString[i+2]
			i += 4
		}
	}
	if toReturnOld {
		return oldValue, nil
	}
	return "", nil
}

func (s *Store) GetExpiryTime(expiryDuration string, divisor int64) (int64, error) {
	millis := time.Now().UnixNano() / divisor
	duration, err := strconv.ParseInt(expiryDuration, 10, 0)
	if err != nil {
		return 0, err
	}
	return millis + duration, nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyMap, ok := s.mp[key]
	if !ok {
		return "", false
	}
	now := time.Now().UnixMilli()
	expires, ok := keyMap["PX"]
	if ok && (now > expires.(int64)) {
		return "", false
	}

	return keyMap["Value"].(string), true
}
