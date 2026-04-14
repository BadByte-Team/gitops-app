package repository

import (
	"database/sql"
	"fmt"

	"github.com/TU_USUARIO/curso-gitops/internal/models"
)

// ── User Repository ──

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db.GetConn()}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password, role, is_blocked, created_at FROM users WHERE username = ?`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.IsBlocked, &user.CreatedAt,
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
	query := `SELECT id, username, role, is_blocked, created_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listando usuarios: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.IsBlocked, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) UpdateRole(id int, role string) error {
	query := `UPDATE users SET role = ? WHERE id = ?`
	_, err := r.db.Exec(query, role, id)
	if err != nil {
		return fmt.Errorf("error actualizando rol: %w", err)
	}
	return nil
}

func (r *UserRepository) ToggleBlock(id int, blocked bool) error {
	query := `UPDATE users SET is_blocked = ? WHERE id = ?`
	_, err := r.db.Exec(query, blocked, id)
	if err != nil {
		return fmt.Errorf("error actualizando bloqueo: %w", err)
	}
	return nil
}

// ── Module Repository ──

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

func (r *ModuleRepository) Create(title string, isHidden bool) (int64, error) {
	query := `INSERT INTO modules (title, is_hidden) VALUES (?, ?)`
	result, err := r.db.Exec(query, title, isHidden)
	if err != nil {
		return 0, fmt.Errorf("error creando módulo: %w", err)
	}
	return result.LastInsertId()
}

func (r *ModuleRepository) Update(id int, title string, isHidden bool) error {
	query := `UPDATE modules SET title = ?, is_hidden = ? WHERE id = ?`
	_, err := r.db.Exec(query, title, isHidden, id)
	if err != nil {
		return fmt.Errorf("error actualizando módulo: %w", err)
	}
	return nil
}

func (r *ModuleRepository) Delete(id int) error {
	query := `DELETE FROM modules WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando módulo: %w", err)
	}
	return nil
}

// ── Episode Repository ──

type EpisodeRepository struct {
	db *sql.DB
}

func NewEpisodeRepository(db *DB) *EpisodeRepository {
	return &EpisodeRepository{db: db.GetConn()}
}

func (r *EpisodeRepository) ListByModule(moduleID int) ([]models.Episode, error) {
	query := `SELECT id, module_id, title, video_url, is_hidden FROM episodes WHERE module_id = ? ORDER BY id`
	rows, err := r.db.Query(query, moduleID)
	if err != nil {
		return nil, fmt.Errorf("error listando episodios: %w", err)
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		var e models.Episode
		if err := rows.Scan(&e.ID, &e.ModuleID, &e.Title, &e.VideoURL, &e.IsHidden); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, nil
}

func (r *EpisodeRepository) ListAll() ([]models.Episode, error) {
	query := `SELECT id, module_id, title, video_url, is_hidden FROM episodes ORDER BY module_id, id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listando episodios: %w", err)
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		var e models.Episode
		if err := rows.Scan(&e.ID, &e.ModuleID, &e.Title, &e.VideoURL, &e.IsHidden); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, nil
}

func (r *EpisodeRepository) Create(moduleID int, title, videoURL string, isHidden bool) (int64, error) {
	query := `INSERT INTO episodes (module_id, title, video_url, is_hidden) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, moduleID, title, videoURL, isHidden)
	if err != nil {
		return 0, fmt.Errorf("error creando episodio: %w", err)
	}
	return result.LastInsertId()
}

func (r *EpisodeRepository) Update(id int, title, videoURL string, isHidden bool) error {
	query := `UPDATE episodes SET title = ?, video_url = ?, is_hidden = ? WHERE id = ?`
	_, err := r.db.Exec(query, title, videoURL, isHidden, id)
	if err != nil {
		return fmt.Errorf("error actualizando episodio: %w", err)
	}
	return nil
}

func (r *EpisodeRepository) Delete(id int) error {
	query := `DELETE FROM episodes WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando episodio: %w", err)
	}
	return nil
}

// ── Progress Repository ──

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *DB) *ProgressRepository {
	return &ProgressRepository{db: db.GetConn()}
}

func (r *ProgressRepository) GetCompletedEpisodes(userID int) ([]int, error) {
	query := `SELECT episode_id FROM user_progress WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error consultando progreso: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *ProgressRepository) ToggleEpisode(userID, episodeID int) (bool, error) {
	// Check if already completed
	var exists int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM user_progress WHERE user_id = ? AND episode_id = ?`, userID, episodeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists > 0 {
		// Remove completion
		_, err = r.db.Exec(`DELETE FROM user_progress WHERE user_id = ? AND episode_id = ?`, userID, episodeID)
		return false, err
	}

	// Add completion
	_, err = r.db.Exec(`INSERT INTO user_progress (user_id, episode_id) VALUES (?, ?)`, userID, episodeID)
	return true, err
}
