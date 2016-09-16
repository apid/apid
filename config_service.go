package apid

import "time"

type ConfigService interface {
	SetDefault(key string, value interface{})
	Set(key string, value interface{})
	Get(key string) interface{}

	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	GetDuration(key string) time.Duration
	IsSet(key string) bool
}
