package services

import (
	"errors"
	"fmt"
	"testovoe/internal/db"
	"testovoe/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	dbProvider *db.PostgresProvider
}

func NewUserService(provider *db.PostgresProvider) *UserService {
	return &UserService{dbProvider: provider}
}

// Хеширование пароля
func (s *UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Проверка пароля
func (s *UserService) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Регистрация нового пользователя
func (s *UserService) RegisterUser(username, password string) (*models.User, error) {
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{Username: username, PasswordHash: hashedPassword}
	err = s.dbProvider.InsertUser(&user)
	return &user, err
}

// Поиск пользователя по логину и проверка пароля
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.dbProvider.GetUserByUsername(username)
	if err != nil {
		fmt.Printf("GetUserByUsername error: %v\n", err)
		return nil, errors.New("invalid username or password")
	}
	if !s.CheckPasswordHash(password, user.PasswordHash) {
		fmt.Println("Password hash mismatch")
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}
