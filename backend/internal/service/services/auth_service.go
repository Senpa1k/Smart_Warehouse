package services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/dgrijalva/jwt-go"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

// AuthService registers users and manages their access tokens.
type AuthService struct {
	repo repository.Authorization
}

// NewAuthService constructs the authentication service.
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser hashes a password and stores the new account.
func (s *AuthService) CreateUser(user models.Users) (uint, error) {
	user.PasswordHash = generateHashPassword(user.PasswordHash)
	return s.repo.CreateUser(user)
}

// GetUser checks credentials and issues a token for a matching user.
func (s *AuthService) GetUser(email, password string) (string, *models.Users, error) {
	// checking the presence of the user in the database
	in, err := s.repo.GetUser(email, generateHashPassword(password))
	if err != nil {
		return "", nil, err
	}

	token := ""
	if in.ID != 0 {
		token = generateJWTToken(in.ID)
	}

	return token, in, nil
}

// ParseToken validates an access token and extracts its user ID.
func (s *AuthService) ParseToken(accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		jwtSecret, _ := config.Get("JWT_SECRET") // jwtSecret будет доступен в контейнере

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not in tokenClaims")
	}

	return claims.UserID, nil
}

func generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	salt, _ := config.Get("SALT") // соль будет доступна в контейнере

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func generateJWTToken(id uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	jwtSecret, _ := config.Get("JWT_SECRET")
	str, _ := token.SignedString([]byte(jwtSecret))
	return str
}
