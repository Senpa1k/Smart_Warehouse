package config

import (
	"fmt"
	"os"
)

func Get(key string) (string, error) {
	// return "postgresql://postgres:arstep2006@localhost:5436/postgres?sslmode=disable", nil
	if val := os.Getenv(key); val != "" {
		return val, nil
	}
	return "", fmt.Errorf("have not acces to env ")
}
