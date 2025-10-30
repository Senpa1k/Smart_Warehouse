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

	// Присваиваем только реальные данные, без пустых записей
	dash.ListScans = scans

	return d.getStatistics(&dash.Statistics)
}

func (d *DashPostgres) getStatistics(statistics *entities.Statistics) error {
	// Общее количество роботов
	var totalRobots int64
	if err := d.db.Model(&models.Robots{}).Count(&totalRobots).Error; err != nil {
		return fmt.Errorf("failed to count total robots: %w", err)
	}
	statistics.TotalRobots = int(totalRobots)

	// Активные роботы
	var activeRobots int64
	if err := d.db.Model(&models.Robots{}).
		Where("status = ?", "active").
		Count(&activeRobots).Error; err != nil {
		return fmt.Errorf("failed to count active robots: %w", err)
	}
	statistics.ActiveRobots = int(activeRobots)

	// Количество проверенных элементов сегодня
	var itemsToday int64
	if err := d.db.Model(&models.InventoryHistory{}).
		Where("DATE(scanned_at) = CURRENT_DATE").
		Count(&itemsToday).Error; err != nil {
		return fmt.Errorf("failed to count items checked today: %w", err)
	}
	statistics.ItemsCheckedToday = int(itemsToday)

	// Критичные элементы (статус CRITICAL)
	var criticalItems int64
	if err := d.db.Model(&models.InventoryHistory{}).
		Where("status = ?", "CRITICAL").
		Count(&criticalItems).Error; err != nil {
		return fmt.Errorf("failed to count critical items: %w", err)
	}
	statistics.CriticalItems = int(criticalItems)

	// Средний заряд батареи роботов
	var avgBattery float64
	row := d.db.Table("robots").
		Select("AVG(battery_level)").
		Where("battery_level IS NOT NULL").
		Row()

	if err := row.Scan(&avgBattery); err != nil {
		statistics.AvgBattery = 0
	} else {
		statistics.AvgBattery = avgBattery
	}

	return nil
}
