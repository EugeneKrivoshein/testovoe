package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testovoe/config"
	"testovoe/internal/db"
	"testovoe/internal/services"
	"testovoe/internal/utils"
	"unicode"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	dbProvider  *db.PostgresProvider
	userService *services.UserService
}

func NewAuthHandler(provider *db.PostgresProvider, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		dbProvider:  provider,
		userService: userService,
	}
}

// Валидация логина и пароля
func (h *AuthHandler) validateCredentials(login, pswd string) error {
	loginRegex := `^[a-zA-Z0-9]{8,}$`

	if matched, _ := regexp.MatchString(loginRegex, login); !matched {
		return errors.New("login must be at least 8 characters long and contain only letters and numbers")
	}

	if len(pswd) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range pswd {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasDigit && hasSpecial) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}

// Регистрация пользователя
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
		Login string `json:"login"`
		Pswd  string `json:"pswd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.Token != config.GetAdminToken() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.validateCredentials(req.Login, req.Pswd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterUser(req.Login, req.Pswd)
	if err != nil {
		fmt.Printf("Error registering user: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]string{"login": user.Username},
	})
}

// Аутентификация пользователя
func (h *AuthHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Login string `json:"login"`
		Pswd  string `json:"pswd"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Printf("Failed to decode request: %v\n", err)
		return
	}

	fmt.Printf("Received login request: login=%s, password=%s\n", req.Login, req.Pswd)

	user, err := h.userService.AuthenticateUser(req.Login, req.Pswd)
	if err != nil {
		fmt.Printf("Authentication error: %v\n", err)
		http.Error(w, "Invalid login or password /", http.StatusUnauthorized)
		return
	}

	token := utils.GenerateToken()
	user.Token = token
	h.dbProvider.UpdateUserToken(user)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]string{"token": token},
	})
}

// Завершение сессии
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Проверьте метод POST
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)    // Извлекаем параметры маршрута
	token := vars["token"] // Получаем токен из параметров маршрута

	if err := h.dbProvider.ClearUserToken(token); err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]bool{token: true},
	})
}
