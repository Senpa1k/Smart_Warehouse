package postgres

import (
	"fmt"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type DashPostgres struct {
	db *gorm.DB
}

func NewDashPostgres(db *gorm.DB) *DashPostgres {
	return &DashPostgres{db: db}
}

func (d *DashPostgres) GetDashInfo(dash *entities.DashInfo) error {
	dash.ListRobots = make([]models.Robots, 0)

	if err := d.db.Order("id").Find(&dash.ListRobots).Error; err != nil {
		return fmt.Errorf("failed to get robots: %w", err)
	}

	var scans []models.InventoryHistory
	if err := d.db.Preload("Robot").Preload("Product").
		Order("scanned_at DESC").
		Limit(100).
		Find(&scans).Error; err != nil {
		return fmt.Errorf("failed to get recent scans: %w", err)
	}

	dash.ListScans = [100]models.InventoryHistory{}

	for i, scan := range scans {
		if i >= 100 {
			break
		}
		dash.ListScans[i] = scan
	}

	return d.getStatistics(&dash.Statistics, dash.ListRobots)
}

func (d *DashPostgres) getStatistics(statistics *entities.Statistics, robots []models.Robots) error {
	// Count active robots
	activeRobots := 0
	totalBattery := 0
	for _, robot := range robots {
		if robot.Status == "active" && robot.ID != "IMPORT_SERVICE" {
			activeRobots++
			totalBattery += robot.BatteryLevel
		}
	}
	
	totalRobots := len(robots)
	if totalRobots > 0 && robots[0].ID == "IMPORT_SERVICE" {
		totalRobots-- // Exclude IMPORT_SERVICE from count
	}
	
	statistics.ActiveRobots = activeRobots
	statistics.TotalRobots = totalRobots
	
	// Average battery
	if activeRobots > 0 {
		statistics.AvgBattery = totalBattery / activeRobots
	} else {
		statistics.AvgBattery = 0
	}
	
	// Items checked today
	today := time.Now().Truncate(24 * time.Hour)
	var itemsCheckedToday int64
	if err := d.db.Model(&models.InventoryHistory{}).
		Where("scanned_at >= ?", today).
		Count(&itemsCheckedToday).Error; err != nil {
		return fmt.Errorf("failed to count items checked today: %w", err)
	}
	statistics.ItemsCheckedToday = int(itemsCheckedToday)
	
	// Critical items (LOW_STOCK or CRITICAL status)
	var criticalItems int64
	if err := d.db.Model(&models.InventoryHistory{}).
		Where("status IN ?", []string{"LOW_STOCK", "CRITICAL"}).
		Distinct("product_id").
		Count(&criticalItems).Error; err != nil {
		return fmt.Errorf("failed to count critical items: %w", err)
	}
	statistics.CriticalItems = int(criticalItems)

	return nil
}
