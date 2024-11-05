package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// Генерация случайного токена
func GenerateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
