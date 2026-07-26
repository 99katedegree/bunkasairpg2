package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := loadEnv(".env"); err != nil {
		slog.Warn("could not load env file", "err", err)
	}

	email := requireEnv("EMAIL")
	password := requireEnv("PASSWORD")

	db, err := connectDB()
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "err", err)
		os.Exit(1)
	}

	_, err = db.ExecContext(context.Background(),
		"INSERT INTO admins (email, password) VALUES (?, ?)",
		email, string(hashed),
	)
	if err != nil {
		slog.Error("failed to create admin", "err", err)
		os.Exit(1)
	}

	slog.Info("admin created", "email", email)
}

// loadEnv は KEY=VALUE 形式の .env ファイルを読んで os.Setenv する。
// 既に環境変数がセットされている場合は上書きしない。
func loadEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
	return scanner.Err()
}

func connectDB() (*sql.DB, error) {
	host := requireEnv("DB_HOST")
	port := envOrDefault("DB_PORT", "3306")
	user := requireEnv("DB_USER")
	pass := requireEnv("DB_PASSWORD")
	name := requireEnv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC", user, pass, host, port, name)
	return sql.Open("mysql", dsn)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "env var %s is required\n", key)
		os.Exit(1)
	}
	return v
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
