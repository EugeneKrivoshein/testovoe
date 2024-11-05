package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"testovoe/config"
	"testovoe/internal/models"

	_ "github.com/lib/pq"
)

type Document struct {
	ID       string
	Name     string
	MimeType string
	File     bool
	Public   bool
	Created  time.Time
	Grant    []string
}

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
