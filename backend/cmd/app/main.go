package main

import (
	"log"

	"github.com/Senpa1k/Smart_Warehouse/internal/handler"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/Senpa1k/Smart_Warehouse/internal/server"
	"github.com/Senpa1k/Smart_Warehouse/internal/service"
)

func main() {
	db, err := repository.InitBD()
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatal("error in init server ", err)
	}
}
