package config

import (
	"os"
	"strconv"
)

func Port() int {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return 0
	}
	return port
}
