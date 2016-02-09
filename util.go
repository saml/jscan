package main

import (
	"os"
	"strconv"
	"strings"
)

func toInt(s string, defaultValue int) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	return defaultValue
}

func envOr(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return defaultValue
}

func firstDir(p string) string {
	return strings.Split(p, "/")[0]
}
