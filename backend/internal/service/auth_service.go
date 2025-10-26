package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.Users) (uint, error) {
	user.PasswordHash = generateHashPassword(user.PasswordHash)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUser(email string, password string) (string, *models.Users, error) {
	// проверка налиия в бд
	in, err := s.repo.GetUser(email, generateHashPassword(password))

	if err != nil {
		return "", nil, err
	}

	token := ""
	if in.ID != 0 {
		token = generateJWTToken()
	}

	return token, in, nil
}

func generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	salt, _ := config.Get("salt")

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func generateJWTToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	jwt_secret, _ := config.Get("jwt_secret")
	str, _ := token.SignedString([]byte(jwt_secret))
	return str
}
