package postgres

import (
	"strings"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type InventoryRepo struct {
	db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) *InventoryRepo {
	return &InventoryRepo{db: db}
}

func (r *InventoryRepo) ImportInventoryHistories(histories []models.InventoryHistory) error {
	return r.db.CreateInBatches(histories, 100).Error
}

func (r *InventoryRepo) GetInventoryHistoryByProductIDs(productIDs []string) ([]models.InventoryHistory, error) {
	var histories []models.InventoryHistory
	err := r.db.Preload("Robot").Preload("Product").Where("product_id IN ?", productIDs).Find(&histories).Error
	return histories, err
}

func (r *InventoryRepo) GetProductByID(productID string) error {
	var product models.Products
	return r.db.First(&product, "id = ?", productID).Error
}

func (r *InventoryRepo) CreateProduct(product *models.Products) error {
	return r.db.Create(product).Error
}

func (r *InventoryRepo) UpdateProduct(product *models.Products) error {
	return r.db.Save(product).Error
}

func parseDateTime(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	
	// Try different date formats
	formats := []string{
		time.RFC3339,           // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano,       // 2006-01-02T15:04:05.999999999Z07:00
		"2006-01-02T15:04:05Z", // ISO 8601 with Z
		"2006-01-02 15:04:05",  // Space separated
		"2006-01-02",           // Date only
	}
	
	var lastErr error
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}
	
	return time.Time{}, lastErr
}

func (r *InventoryRepo) GetHistory(from, to, zone, status string, limit, offset int) ([]models.InventoryHistory, int64, error) {
	var histories []models.InventoryHistory
	var total int64

	query := r.db.Model(&models.InventoryHistory{})

	if from != "" {
		filterTime, err := parseDateTime(from)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("scanned_at >= ?", filterTime)
	}

	if to != "" {
		filterTime, err := parseDateTime(to)
		if err != nil {
			return nil, 0, err
		}
		// If time is exactly midnight (no time component), add full day
		if filterTime.Hour() == 0 && filterTime.Minute() == 0 && filterTime.Second() == 0 {
			filterTime = filterTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		query = query.Where("scanned_at <= ?", filterTime)
	}

	if zone != "" {
		query = query.Where("zone = ?", zone)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total before applying limit/offset
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply preload, limit, offset and order
	err := query.Preload("Robot").Preload("Product").Limit(limit).Offset(offset).Order("scanned_at DESC").Find(&histories).Error
	return histories, total, err
}
