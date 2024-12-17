package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresConnection() (*sql.DB, error) {
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnvAsInt("DB_PORT", 5432)
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "orders")

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var db *sql.DB
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Tentativa %d: Erro ao abrir conexão com o banco: %v", i+1, err)
			time.Sleep(time.Second * 2)
			continue
		}

		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Tentativa %d: Erro ao pingar o banco: %v", i+1, err)
		db.Close()
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados após %d tentativas: %v", maxRetries, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
