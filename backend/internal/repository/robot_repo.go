package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type RobotPostgres struct {
	db *gorm.DB
}

func NewRobotPostgres(db *gorm.DB) *RobotPostgres {
	return &RobotPostgres{db: db}
}

func (r *RobotPostgres) AddData(data entities.RobotsData) error { // обработать ошибки типа неправ знач в поле
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, scanResult := range data.ScanResults {
		//проверка foreignkey и id_robot
		var count int64
		if tx.Model(&models.Products{}).Where("id = ?", scanResult.ProductId).Count(&count); count == 0 {
			tx.Rollback()
			return fmt.Errorf("product does not exist")
		}
		if tx.Model(&models.Robots{}).Where("id = ?", data.RobotId).Count(&count); count == 0 {
			tx.Rollback()
			return fmt.Errorf("robot does not exist")
		}

		var inventoryHistory models.InventoryHistory = models.InventoryHistory{
			RobotID:     data.RobotId,
			ProductID:   scanResult.ProductId,
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

	nextPoint := strings.Split(data.NextCheckpoint, "-")
	row, err1 := strconv.Atoi(nextPoint[1])
	shelf, err2 := strconv.Atoi(nextPoint[2])
	if err1 != nil {
		tx.Rollback()
		return err1
	}
	if err2 != nil {
		tx.Rollback()
		return err2
	}

	var robot models.Robots
	if err := tx.Where("id = ?", data.RobotId).First(&robot).Error; err != nil {
		tx.Rollback()
		return err
	}
	robot.BatteryLevel = data.BatteryLevel
	robot.CurrentZone = nextPoint[0]
	robot.CurrentRow = row
	robot.CurrentShelf = shelf
	if err := tx.Save(&robot).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *RobotPostgres) CheckId(robotID string) bool {
	var count int64
	err := r.db.Model(&models.Robots{}).Where("id = ?", robotID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
