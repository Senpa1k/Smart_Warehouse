package repository

import (
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

func (r *RobotPostgres) AddData(data entities.RobotsData) error {
	for _, scanResult := range data.ScanResults {
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

		if err := r.db.Create(&inventoryHistory).Error; err != nil {
			return err
		}
	}

	nextPoint := strings.Split(data.NextCheckpoint, "-")
	row, err1 := strconv.Atoi(nextPoint[1])
	shelf, err2 := strconv.Atoi(nextPoint[2])
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	var robot models.Robots
	r.db.First(&robot, data.RobotId)
	robot.BatteryLevel = data.BatteryLevel
	robot.CurrentZone = nextPoint[0]
	robot.CurrentRow = row
	robot.CurrentShelf = shelf
	r.db.Save(&robot)

	return nil
}
