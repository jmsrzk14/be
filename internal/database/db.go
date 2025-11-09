package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"os/exec"

	"bem_be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database connection
var DB *gorm.DB

// Initialize connects to the Supabase PostgreSQL database and creates tables if they don't exist
func Initialize() {
	var err error

	// Ambil konfigurasi dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Pastikan semua variabel environment sudah ada
	if host == "" || port == "" || user == "" || password == "" || dbname == "" || sslmode == "" {
		log.Fatal("Salah satu environment variable database belum diatur")
	}

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	// Buat DSN PostgreSQL untuk Supabase
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	// Connect ke database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Error connecting to Supabase database: %v", err)
	}

	// Atur connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying SQL DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to Supabase PostgreSQL database successfully")

	// Migration sequence
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Category{},
		&models.Organization{},
		&models.Period{},
		&models.Student{},
		&models.BEM{},
		&models.Activity{},
		&models.News{},
		&models.Galery{},
		&models.Item{},
		&models.Request{},
		&models.Aspiration{},
		&models.Calender{},
		&models.Announcement{},
		&models.MPM{},
		&models.Notification{},
		&models.UserNotification{},
		&models.Status_Aspirations{},
	}

	for _, model := range modelsToMigrate {
		err = DB.AutoMigrate(model)
		if err != nil {
			log.Fatalf("Error auto-migrating model %T: %v", model, err)
		}
		log.Printf("%T table migrated successfully", model)
	}

	log.Println("Database schema migrated successfully")
}

// Close closes the database connection
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error getting underlying SQL DB: %v", err)
			return
		}
		sqlDB.Close()
	}
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}

func BackupDatabase() {
	// Ambil konfigurasi dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Println("Environment variable database belum lengkap. Backup dibatalkan")
		return
	}

	// Buat folder backup jika belum ada
	backupDir := "backups"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err = os.Mkdir(backupDir, 0755)
		if err != nil {
			log.Fatalf("Gagal membuat folder backup: %v", err)
		}
	}

	// Format nama file backup: backup-YYYY-MM.sql
	now := time.Now()
	fileName := fmt.Sprintf("%s/backup-%d.sql", backupDir, now.Year())

	// Set environment variable PGPASSWORD supaya pg_dump bisa otomatis login
	os.Setenv("PGPASSWORD", password)

	// Jalankan pg_dump
	cmd := exec.Command(
		"pg_dump",
		"-h", host,
		"-p", port,
		"-U", user,
		"-F", "p", // format custom
		"-a",      // include large objects
		"-f", fileName,
		dbname,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Gagal melakukan backup database: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Backup database berhasil: %s\n", fileName)
}