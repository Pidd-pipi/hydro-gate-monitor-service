package config

import (
	"os"
	"strconv"
)

func Port() int {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port < 1 || port > 65535 {
		return 8080
	}
	return port
}
