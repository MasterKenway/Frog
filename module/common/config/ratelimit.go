package config

type RateLimitConfig struct {
	TimesPerSec uint64 `json:"times_per_sec"`
}
