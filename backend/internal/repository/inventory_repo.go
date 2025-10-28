package repository

import (
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

func (r *InventoryRepo) GetHistory(from, to, zone, status string, limit, offset int) ([]models.InventoryHistory, int64, error) {
	var histories []models.InventoryHistory
	var total int64

	query := r.db.Model(&models.InventoryHistory{}).Preload("Robot").Preload("Product")

	if from != "" {
		filterTime, err := time.Parse("2006-01-02 15:04:05", from)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("scanned_at >= ?", filterTime)
	}

	if to != "" {
		filterTime, err := time.Parse("2006-01-02 15:04:05", to)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("scanned_at <= ?", filterTime)
	}

	if zone != "" {
		query = query.Where("zone = ?", zone)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Limit(limit).Offset(offset).Order("scanned_at DESC").Find(&histories).Error
	return histories, total, err
}
