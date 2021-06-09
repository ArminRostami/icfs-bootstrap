package env

import (
	"os"
	"strings"
)

const LocalNetIP = "192.168.56.1"

func DockerEnabled() bool {
	val, exists := os.LookupEnv("DOCKER_ENABLED")
	if !exists {
		return false
	}
	return strings.EqualFold(val, "1")
}
