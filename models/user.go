package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User struct represents a user in our system
type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	LastLogin    *time.Time
}

// User Repository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database
func (r *UserRepository) CreateUser(email, username, passwordHash string) (*User, error) {
	user := &User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO users (id, email, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retreives a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, last_login 
		FROM users
		WHERE email = $1;
	`
	var user User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &lastLogin)
	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

// GetUserByID retreives a user by their ID
func (r *UserRepository) GetUserByID(id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, last_login
		FROM users
		WHERE id = $1;
	`

	var user User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&lastLogin)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}
