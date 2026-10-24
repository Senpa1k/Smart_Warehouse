package postgres

import (
	"fmt"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

// RobotPostgres persists robot scans and positions in PostgreSQL.
type RobotPostgres struct {
	db *gorm.DB
}

// NewRobotPostgres constructs the PostgreSQL robot repository.
func NewRobotPostgres(db *gorm.DB) *RobotPostgres {
	return &RobotPostgres{db: db}
}

// AddData stores a robot scan and updates its current position.
func (r *RobotPostgres) AddData(data entities.RobotsData) error { // обработать ошибки типа неправ знач в поле
	zone, row, shelf, err := entities.ParseCheckpoint(data.NextCheckpoint)
	if err != nil {
		return err
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// обработка результатов сканирования роботов
	for _, scanResult := range data.ScanResults {
		// проверка foreignkey и id_robot
		var count int64
		if tx.Model(&models.Products{}).Where("id = ?", scanResult.ProductID).Count(&count); count == 0 {
			tx.Rollback()
			return fmt.Errorf("product does not exist")
		}
		if tx.Model(&models.Robots{}).Where("id = ?", data.RobotID).Count(&count); count == 0 {
			tx.Rollback()
			return fmt.Errorf("robot does not exist")
		}

		// построение экземпляра структуры истории инвентаризации
		inventoryHistory := models.InventoryHistory{
			RobotID:     data.RobotID,
			ProductID:   scanResult.ProductID,
			Quantity:    scanResult.Quantity,
			Zone:        data.Location.Zone,
			RowNumber:   data.Location.Row,
			ShelfNumber: data.Location.Shelf,
			Status:      scanResult.Status,
			ScannedAt:   data.Timestamp,
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&inventoryHistory).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// парсинг информации о роботе
	// получение робота по id и обновление его полей
	var robot models.Robots
	if err := tx.Where("id = ?", data.RobotID).First(&robot).Error; err != nil {
		tx.Rollback()
		return err
	}
	robot.BatteryLevel = data.BatteryLevel
	robot.CurrentZone = zone
	robot.CurrentRow = row
	robot.CurrentShelf = shelf
	if err := tx.Save(&robot).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CheckID reports whether a robot with the given ID exists.
func (r *RobotPostgres) CheckID(robotID string) bool {
	var count int64
	err := r.db.Model(&models.Robots{}).Where("id = ?", robotID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
