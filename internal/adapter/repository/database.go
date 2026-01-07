package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Подключается к бд
func newDatabaseConn(cfg DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Database connected!")

	return db, nil
}

// Проверяет наличие базы данных и при ее отсуствии создает
func createDatabaseIfNotExists(cfg DatabaseConfig) error {
	tempCfg := cfg
	tempCfg.DBName = "postgres"

	db, err := newDatabaseConn(tempCfg)
	if err != nil {
		return fmt.Errorf("ошибка при подключении к системной БД: %w", err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE DATNAME = $1)",
		cfg.DBName,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования БД: %w", err)
	}

	if !exists {
		if _, err := db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)); err != nil {
			return fmt.Errorf("ошибка создания БД: %w", err)
		}
		log.Printf("✅ База данных '%s' создана\n", cfg.DBName)
	} else {
		log.Printf("✅ База данных '%s', уже существует\n", cfg.DBName)
	}

	return nil
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil{
		return fmt.Errorf("ошибка чтения папки миграций: %w", err)
	}

	var sqlFiles []string
	for _, file := range files{
		if filepath.Ext(file.Name()) == ".sql"{
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Println("⚠️ Миграции не найдены")
		return nil
	}

	for _, fileName := range sqlFiles{
		filePath := filepath.Join(migrationsDir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil{
			return fmt.Errorf("ошибка чтения файла %s: %w", fileName, err)
		}

		if _, err := db.Exec(string(content)); err != nil{
			return fmt.Errorf("ошибка выполнения %s: %w", fileName, err)
		}

		log.Printf("✅ Миграция %s выполнена\n", fileName)
	}

	return nil
}

// InitDatabase инициализирует БД: создаёт её, подключается и запускает миграции
func InitDatabase(cfg DatabaseConfig, migrationDir string) (*sql.DB, error) {
	if err := createDatabaseIfNotExists(cfg); err != nil {
		return nil, err
	}

	db, err := newDatabaseConn(cfg)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db, migrationDir); err != nil {
		return nil, err
	}

	log.Println("✅ БД полностью инициализирована!")
	return db, nil
}

// CloseDB закрывает соединение с БД
func CloseDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия БД: %w", err)
	}
	log.Println("✅ Соединение с БД закрыто")
	return nil
}