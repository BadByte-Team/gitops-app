package repository

import (
	"database/sql"
	"fmt"

	"github.com/TU_USUARIO/curso-gitops/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db.GetConn()}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password, role, created_at FROM users WHERE username = ?`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("usuario no encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando usuario: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Create(username, hashedPassword, role string) error {
	query := `INSERT INTO users (username, password, role) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, username, hashedPassword, role)
	if err != nil {
		return fmt.Errorf("error creando usuario: %w", err)
	}
	return nil
}

func (r *UserRepository) ListAll() ([]models.User, error) {
	query := `SELECT id, username, role, created_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listando usuarios: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

type ModuleRepository struct {
	db *sql.DB
}

func NewModuleRepository(db *DB) *ModuleRepository {
	return &ModuleRepository{db: db.GetConn()}
}

func (r *ModuleRepository) ListVisible() ([]models.Module, error) {
	query := `SELECT id, title, is_hidden FROM modules WHERE is_hidden = FALSE ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listando módulos: %w", err)
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var m models.Module
		if err := rows.Scan(&m.ID, &m.Title, &m.IsHidden); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *ModuleRepository) ListAll() ([]models.Module, error) {
	query := `SELECT id, title, is_hidden FROM modules ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listando módulos: %w", err)
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var m models.Module
		if err := rows.Scan(&m.ID, &m.Title, &m.IsHidden); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}
