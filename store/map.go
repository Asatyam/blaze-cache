package store

import (
	"strings"
	"sync"
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

func (s *Store) Set(arrString []string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// arr would be of form ["sym", "key", "sym", "value", ...opts  ]
	//opts ["sym", "opt", "sym", "optValue",..........]

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
		if strings.ToUpper(arrString[i]) == "GET" {
			toReturnOld = true
			i += 2
		} else {
			s.mp[key][arrString[i]] = arrString[i+2]
			i += 4
		}
	}
	if toReturnOld {
		return oldValue, true, nil
	}
	return "", true, nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyMap, ok := s.mp[key]
	if !ok {
		return "", false
	}
	return keyMap["Value"].(string), true
}
