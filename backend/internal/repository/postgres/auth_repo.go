package postgres

import (
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// adding a user to the database
func (r *AuthPostgres) CreateUser(user models.Users) (uint, error) {
	result := r.db.Create(&user)
	if err := result.Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// getting user data from the database by password and email
func (r *AuthPostgres) GetUser(email string, passwordHash string) (*models.Users, error) {
	var user models.Users
	result := r.db.Where("password_hash = ? and email = ?", passwordHash, email).Find(&user)
	return &user, result.Error
}
