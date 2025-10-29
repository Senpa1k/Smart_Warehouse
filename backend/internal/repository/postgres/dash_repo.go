package postgres

import (
	"fmt"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type DashPostgres struct {
	db *gorm.DB
}

type Statistics struct {
	TotalCheck        int `json:"total_check"`
	UniqueProducts    int `json:"unique_products"`
	FindDiscrepancies int `json:"find_discrepancies"`
	AverageTime       int `json:"average_time"`
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

	return d.getStatistics(&dash.Statistics)
}

func (d *DashPostgres) getStatistics(statistics *entities.Statistics) error {
	// Общее количество проверок
	var totalCheck int64
	if err := d.db.Model(&models.InventoryHistory{}).Count(&totalCheck).Error; err != nil {
		return fmt.Errorf("failed to count total checks: %w", err)
	}
	statistics.TotalCheck = int(totalCheck)

	var uniqueProducts int64
	if err := d.db.Model(&models.InventoryHistory{}).
		Distinct("product_id").
		Count(&uniqueProducts).Error; err != nil {
		return fmt.Errorf("failed to count unique products: %w", err)
	}
	statistics.UniqueProducts = int(uniqueProducts)

	var discrepancies int64 // расхождения?
	// if err := d.db.Model(&models.InventoryHistory{}).
	// 	Where("status = ?", "discrepancy").
	// 	Count(&discrepancies).Error; err != nil {
	// 	return fmt.Errorf("failed to count discrepancies: %w", err)
	// }
	statistics.FindDiscrepancies = int(discrepancies)

	var avgTime float64
	row := d.db.Table("inventory_histories").
		Select("AVG(EXTRACT(EPOCH FROM (scanned_at - created_at)))").
		Where("scanned_at > created_at").
		Row()

	if err := row.Scan(&avgTime); err != nil {
		statistics.AverageTime = 0
	} else {
		statistics.AverageTime = int(avgTime)
	}

	return nil
}
