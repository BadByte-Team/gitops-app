package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
}

type Module struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	IsHidden bool   `json:"is_hidden"`
}

type Episode struct {
	ID       int    `json:"id"`
	ModuleID int    `json:"module_id"`
	Title    string `json:"title"`
	VideoURL string `json:"video_url"`
	IsHidden bool   `json:"is_hidden"`
}

type UserProgress struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	EpisodeID   int       `json:"episode_id"`
	CompletedAt time.Time `json:"completed_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ModuleRequest struct {
	Title    string `json:"title"`
	IsHidden bool   `json:"is_hidden"`
}

type EpisodeRequest struct {
	ModuleID int    `json:"module_id"`
	Title    string `json:"title"`
	VideoURL string `json:"video_url"`
	IsHidden bool   `json:"is_hidden"`
}

type UserRoleRequest struct {
	Role string `json:"role"`
}

type ToggleProgressRequest struct {
	EpisodeID int `json:"episode_id"`
}
