package postgres

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

// AuthPostgres stores and looks up user accounts in PostgreSQL.
type AuthPostgres struct {
	db *gorm.DB
}

// NewAuthPostgres constructs the PostgreSQL authentication repository.
func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// CreateUser persists a user account and returns its ID.
func (r *AuthPostgres) CreateUser(user models.Users) (uint, error) {
	result := r.db.Create(&user)
	if err := result.Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// GetUser finds a user by email and password hash.
func (r *AuthPostgres) GetUser(email, passwordHash string) (*models.Users, error) {
	var user models.Users
	result := r.db.Where("password_hash = ? and email = ?", passwordHash, email).Find(&user)
	return &user, result.Error
}
