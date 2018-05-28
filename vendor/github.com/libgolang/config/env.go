package config

import (
	"os"
	"strconv"
	"strings"
)

func envString(name string, def string) string {
	if v, ok := os.LookupEnv(_env(name)); ok {
		return v
	}
	return getConfig().GetString(name, def)
}

func envInt(name string, def int) int {
	if v, ok := os.LookupEnv(_env(name)); ok {
		i, _ := strconv.Atoi(v)
		return i
	}
	return getConfig().GetInt(name, def)
}

func envInt64(name string, def int64) int64 {
	if v, ok := os.LookupEnv(_env(name)); ok {
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	}
	return getConfig().GetInt64(name, def)
}

func envFloat(name string, def float64) float64 {
	if v, ok := os.LookupEnv(_env(name)); ok {
		i, _ := strconv.ParseFloat(v, 64)
		return i
	}
	return getConfig().GetFloat64(name, def)
}

func envBool(name string, def bool) bool {
	if v, ok := os.LookupEnv(_env(name)); ok {
		i, _ := strconv.ParseBool(v)
		return i
	}
	return getConfig().GetBool(name, def)
}

func _env(key string) string {
	envKey := strings.Replace(key, ".", "_", -1)
	envKey = strings.Replace(key, "-", "_", -1)
	envKey = strings.ToUpper(key)
	return envKey
}
