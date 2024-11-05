package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"testovoe/config"
	"testovoe/internal/models"

	"github.com/lib/pq"
)

type PostgresProvider struct {
	db *sql.DB
}

func NewPostgresProvider() (*PostgresProvider, error) {
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database.")

	return &PostgresProvider{db: db}, nil
}
func (p *PostgresProvider) InsertUser(user *models.User) error {
	_, err := p.db.Exec(
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)`,
		user.Username, user.PasswordHash,
	)

	if err != nil {
		fmt.Printf("InsertUser error: %v\n", err)
		return err
	}

	return err
}

func (p *PostgresProvider) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	var token sql.NullString

	err := p.db.QueryRow(
		`SELECT id, username, password_hash, token FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Token)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if token.Valid {
		user.Token = token.String
	} else {
		user.Token = ""
	}

	return user, err
}

func (p *PostgresProvider) GetUserByToken(token string) (*models.User, error) {
	user := &models.User{}
	err := p.db.QueryRow(
		`SELECT id, username, password_hash, token FROM users WHERE token = $1`,
		token,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Token)

	if err == sql.ErrNoRows {
		return nil, nil // пользователь с таким токеном не найден
	}
	if err != nil {
		return nil, err // ошибка базы данных
	}

	return user, nil
}

func (p *PostgresProvider) UpdateUserToken(user *models.User) error {
	_, err := p.db.Exec(
		`UPDATE users SET token = $1 WHERE id = $2`,
		user.Token, user.ID,
	)
	return err
}

func (p *PostgresProvider) ClearUserToken(token string) error {
	_, err := p.db.Exec(`UPDATE users SET token = NULL WHERE token = $1`, token)
	if err != nil {
		return errors.New("failed to clear token")
	}
	return nil
}

func (p *PostgresProvider) Close() error {
	return p.db.Close()
}

func (p *PostgresProvider) SaveDocument(doc *models.Document) error {
	_, err := p.db.Exec(
		`INSERT INTO documents (user_id, name, mime_type, content, public, grant) VALUES ($1, $2, $3, $4, $5, $6)`,
		doc.UserID, doc.Name, doc.MimeType, doc.Content, doc.Public, pq.Array(doc.Grant),
	)
	return err
}
