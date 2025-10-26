package main

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/handler"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/Senpa1k/Smart_Warehouse/internal/server"
	"github.com/Senpa1k/Smart_Warehouse/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	db, err := repository.InitBD()
	if err != nil {
		logrus.Fatalf("fatal initializetion db, %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		logrus.Fatalf("error in init http server: %s", err.Error())
	}
}
