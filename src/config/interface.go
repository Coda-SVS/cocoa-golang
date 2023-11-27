package config

import (
	"strings"
	"time"

	"github.com/spf13/cast"
)

func (c *Config) Set(key string, value any) {
	if c.currentPath == "" {
		c.v.Set(key, value)
	} else {
		c.v.Set(strings.Join([]string{c.currentPath, key}, "."), value)
	}
}

func (c *Config) SetDefault(key string, value any) {
	if c.currentPath == "" {
		c.v.SetDefault(key, value)
	} else {
		c.v.SetDefault(strings.Join([]string{c.currentPath, key}, "."), value)
	}
}

func (c *Config) GetSub(path string) *Config {
	cfg := &Config{v: c.v}
	if c.currentPath == "" {
		cfg.currentPath = path
	} else {
		cfg.currentPath = strings.Join([]string{c.currentPath, path}, ".")
	}
	return cfg
}

func (c *Config) Get(key string) any {
	if c.currentPath == "" {
		return c.v.Get(key)
	} else {
		return c.v.Get(strings.Join([]string{c.currentPath, key}, "."))
	}
}

func (c *Config) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

func (c *Config) GetInt32(key string) int32 {
	return cast.ToInt32(c.Get(key))
}

func (c *Config) GetInt64(key string) int64 {
	return cast.ToInt64(c.Get(key))
}

func (c *Config) GetUint(key string) uint {
	return cast.ToUint(c.Get(key))
}

func (c *Config) GetUint16(key string) uint16 {
	return cast.ToUint16(c.Get(key))
}

func (c *Config) GetUint32(key string) uint32 {
	return cast.ToUint32(c.Get(key))
}

func (c *Config) GetUint64(key string) uint64 {
	return cast.ToUint64(c.Get(key))
}

func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

func (c *Config) GetTime(key string) time.Time {
	return cast.ToTime(c.Get(key))
}

func (c *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}

func (c *Config) GetIntSlice(key string) []int {
	return cast.ToIntSlice(c.Get(key))
}

func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}

func (c *Config) GetStringMap(key string) map[string]any {
	return cast.ToStringMap(c.Get(key))
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.Get(key))
}

func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(c.Get(key))
}
