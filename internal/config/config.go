package config

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Config struct {
	mu sync.RWMutex
	mp map[string]string
}

func NewConfig(dir, dbFileName string) *Config {
	mp := make(map[string]string)
	mp["dir"] = dir
	mp["dbfilename"] = dbFileName
	return &Config{
		mp: mp,
	}
}

func (c *Config) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	key = strings.ToLower(key)
	value, ok := c.mp[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}
