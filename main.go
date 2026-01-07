package main

import (
	"fuelStation/internal/adapter/repository"
	"log"
)

func main() {
	// Подключение
	db, err := repository.InitDatabase(repository.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123qwe",
		DBName:   "gasStation",
		SSLMode:  "disable",
	}, "internal/adapter/repository/migrations")
	if err != nil {
		log.Fatal("❌", err)
	}
	defer repository.CloseDB(db)

	// Проверка
	err = db.Ping()
	if err != nil {
		log.Fatal("❌", err)
	}

	log.Println("✅ Database OK!")
}
