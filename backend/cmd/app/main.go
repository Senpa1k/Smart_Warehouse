package main

import (
	"context"

	"github.com/Senpa1k/Smart_Warehouse/internal/delivery/http/handler"
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

	done := make(chan struct{})

	srv := new(server.Server)
	go func() {
		if err := srv.Run("8080", handler.InitRoutes()); err != nil {
			logrus.Fatalf("error in init http server: %s", err.Error())
		}
		done <- struct{}{}
	}()

	logrus.Print("Server up")
	<-done
	logrus.Print("Server down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error with shutting down %s", err.Error())
	}

	closer, err := db.DB()
	if err2 := closer.Close(); err != nil || err2 != nil {
		logrus.Fatalf("error with clossing db %s", err.Error())
	}
}
