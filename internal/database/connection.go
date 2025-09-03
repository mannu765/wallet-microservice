package database

import (
	"fmt"
	"log"
	"os"
	"wallet-microservice/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "wallet_db")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func Migrate() {
	// Enable UUID extension
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	// Auto-migrate the models with all GORM tags
	err := DB.AutoMigrate(&models.Wallet{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create updated_at trigger for wallets table
	createUpdatedAtTrigger()

	log.Println("Database migration completed")
}

func createUpdatedAtTrigger() {
	// Create the trigger function if it doesn't exist
	triggerFunction := `
        CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = CURRENT_TIMESTAMP;
            RETURN NEW;
        END;
        $$ language 'plpgsql';
    `

	if err := DB.Exec(triggerFunction).Error; err != nil {
		log.Printf("Warning: Failed to create trigger function: %v", err)
		return
	}

	// Create the trigger if it doesn't exist
	trigger := `
        DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;
        CREATE TRIGGER update_wallets_updated_at 
        BEFORE UPDATE ON wallets
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    `

	if err := DB.Exec(trigger).Error; err != nil {
		log.Printf("Warning: Failed to create trigger: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
