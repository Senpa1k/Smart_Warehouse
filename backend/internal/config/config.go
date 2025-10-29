package config

import (
	"fmt"
	"os"
)

func Get(key string) (string, error) {
	if val := os.Getenv(key); val != "" {
		return val, nil
	}
	return "", fmt.Errorf("have not acces to env ")
}
