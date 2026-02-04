package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type AgentConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	HostProc      string
	Interval      time.Duration
	ChannelPrefix string
	DiskPath      string
}

type HubConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Pattern       string
	WSPath        string
	HTTPSAddr     string
	TLSCertPath   string
	TLSKeyPath    string
	StaticDir     string
	Namespace     string
	ScaleTargets  []string
	MaxReplicas   int32
	ScaleEnabled  bool
}

type TraceSvcConfig struct {
	HTTPAddr string
}

func LoadAgentConfig() AgentConfig {
	return AgentConfig{
		RedisAddr:     getenv("REDIS_ADDR", "redis-service:6379"),
		RedisPassword: getenv("REDIS_PASSWORD", ""),
		RedisDB:       getenvInt("REDIS_DB", 0),
		HostProc:      getenv("HOST_PROC", "/host/proc"),
		Interval:      getenvDuration("SCRAPE_INTERVAL", 1*time.Second),
		ChannelPrefix: getenv("CHANNEL_PREFIX", "metrics"),
		DiskPath:      getenv("DISK_PATH", "/"),
	}
}

func LoadHubConfig() HubConfig {
	return HubConfig{
		RedisAddr:     getenv("REDIS_ADDR", "redis-service:6379"),
		RedisPassword: getenv("REDIS_PASSWORD", ""),
		RedisDB:       getenvInt("REDIS_DB", 0),
		Pattern:       getenv("REDIS_PATTERN", "metrics:*"),
		WSPath:        getenv("WS_PATH", "/ws"),
		HTTPSAddr:     getenv("HTTPS_ADDR", ":8443"),
		TLSCertPath:   getenv("TLS_CERT_PATH", "/tls/tls.crt"),
		TLSKeyPath:    getenv("TLS_KEY_PATH", "/tls/tls.key"),
		StaticDir:     getenv("STATIC_DIR", "web"),
		Namespace:     getenv("NAMESPACE", "default"),
		ScaleTargets:  splitList(getenv("SCALE_TARGETS", "resource-hub,redis")),
		MaxReplicas:   int32(getenvInt("SCALE_MAX", 5)),
		ScaleEnabled:  getenvBool("SCALE_ENABLED", true),
	}
}

func LoadTraceSvcConfig() TraceSvcConfig {
	return TraceSvcConfig{
		HTTPAddr: getenv("HTTP_ADDR", ":8081"),
	}
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getenvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return parsed
}

func getenvDuration(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return parsed
}

func getenvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return parsed
}

func splitList(val string) []string {
	if val == "" {
		return nil
	}
	raw := make([]string, 0)
	for _, item := range strings.Split(val, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		raw = append(raw, trimmed)
	}
	return raw
}
