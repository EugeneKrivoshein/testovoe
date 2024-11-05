package models

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Token        string `json:"token,omitempty"`
}

type Document struct {
	ID       int
	UserID   int
	Name     string
	MimeType string
	File     bool
	Public   bool
	Grant    []string
	Content  []byte
	Created  time.Time
}
