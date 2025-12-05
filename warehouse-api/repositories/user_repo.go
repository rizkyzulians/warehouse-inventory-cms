package repositories

import (
	"database/sql"
	"warehouse-api/models"
)

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByID(id int) (*models.User, error)
	Create(user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password, nama, role, created_at, updated_at 
	          FROM users WHERE username = $1`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Nama,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password, nama, role, created_at, updated_at 
	          FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Nama,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, password, nama, role) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, user.Username, user.Password, user.Nama, user.Role).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
}
