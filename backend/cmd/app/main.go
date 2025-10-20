package main

import (
	"log"

	"github.com/Senpa1k/Smart_Warehouse/internal/config"
)

func main() {
	_, err := config.InitDB()
	if err != nil {
		log.Fatal("can not open db")
	}

}
