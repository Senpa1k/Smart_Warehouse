package main

import (
	"log"
	"smart_warehouse/backend/internal/config"
)

func main() {
	_, err := config.InitDB()
	if err != nil {
		log.Fatal("can not open db")
	}

}
