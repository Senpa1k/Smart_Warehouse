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

	// Инициализация Redis
	redisClient, err := repository.NewRedisClient("redis://localhost:6379")
	if err != nil {
		logrus.Warnf("Redis connection failed: %v", err)
		logrus.Info("Application will continue without Redis caching")
		redisClient = nil
	} else {
		logrus.Info("Redis connected successfully")
		defer redisClient.Close()
	}

	// Создаем репозитории с Redis
	repos := repository.NewRepository(db, redisClient)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)

	done := make(chan struct{})

	srv := new(server.Server)
	go func() {
		if err := srv.Run("3000", handler.InitRoutes()); err != nil {
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
