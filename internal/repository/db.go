package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "curso_app")
	pass := getEnv("DB_PASSWORD", "C4rs0_S3cur3_P@ss!")
	name := getEnv("DB_NAME", "curso_db")
	port := getEnv("DB_PORT", "3306")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexión: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a la BD: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return &DB{conn: db}, nil
}

func (d *DB) GetConn() *sql.DB {
	return d.conn
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
