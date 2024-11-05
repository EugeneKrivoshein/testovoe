// handlers/doc_handler.go
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testovoe/internal/db"
	"testovoe/internal/middleware"
	"testovoe/internal/models"
	"testovoe/internal/services"
)

type DocumentHandler struct {
	dbProvider *db.PostgresProvider
	DocService *services.DocumentService
}

func NewDocumentHandler(provider *db.PostgresProvider, docService *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		dbProvider: provider,
		DocService: docService,
	}
}

func (h *DocumentHandler) UploadDocumentHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Ограничение 10 MB для файла
	if err != nil {
		http.Error(w, "Ошибка при обработке формы", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Чтение метаданных из параметра "meta"
	var meta struct {
		Name   string   `json:"name"`
		File   bool     `json:"file"`
		Public bool     `json:"public"`
		Token  string   `json:"token"`
		Mime   string   `json:"mime"`
		Grant  []string `json:"grant"`
	}
	err = json.Unmarshal([]byte(r.FormValue("meta")), &meta)
	if err != nil {
		http.Error(w, "Ошибка при разборе метаданных", http.StatusBadRequest)
		return
	}

	// Считывание файла в []byte для хранения в базе
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	// Создаем документ
	doc := models.Document{
		Name:     user.Username,
		MimeType: meta.Mime,
		File:     meta.File,
		Public:   meta.Public,
		Grant:    meta.Grant,
		Content:  fileBytes,
	}

	// Сохраняем документ через сервис
	if err := h.dbProvider.SaveDocument(&doc); err != nil {
		http.Error(w, "Ошибка при сохранении документа", http.StatusInternalServerError)
		return
	}

	// Ответ с подтверждением
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]interface{}{
			"json": doc.Content,
			"file": handler.Filename,
		},
	})
}
