/*
Код для принятия миграций из migrates. Запускается автоматически при сборке контейнера backend
*/

// Package main applies warehouse database migrations.
package main

import (
	"errors"
	"log"

	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func runMigrate(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Println("Migrations applied successfully")
}

func main() {
	db, err := repository.InitBD()
	if err != nil {
		log.Fatal(err)
	}

	runMigrate(db)

	log.Print("everything is good")
}
