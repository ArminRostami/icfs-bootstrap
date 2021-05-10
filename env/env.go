package env

import (
	"os"
	"strings"
)

func DockerEnabled() bool {
	val, exists := os.LookupEnv("DOCKER_ENABLED")
	if !exists {
		return false
	}
	if strings.EqualFold(val, "1") {
		return true
	}
	return false
}
