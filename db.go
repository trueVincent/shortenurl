package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedOn time.Time `gorm:"autoCreateTime"`
	UpdatedOn time.Time `gorm:"autoUpdateTime"`
}

type URLMapping struct {
	ID        string `gorm:"size:6;uniqueIndex;not null"`
	OriginURL string `gorm:"not null"`
	UserID    uint
	User      User
	CreatedOn time.Time `gorm:"autoCreateTime"`
	UpdatedOn time.Time `gorm:"autoUpdateTime"`
}

type URLMappingActionRecord struct {
	ID           uint   `gorm:"primaryKey"`
	URLMappingID string `gorm:"not null"`
	URLMapping   URLMapping
	ClickCount   int
	LastAccess   time.Time `gorm:"autoUpdateTime"`
	CreatedOn    time.Time `gorm:"autoCreateTime"`
	UpdatedOn    time.Time `gorm:"autoUpdateTime"`
}

func InitializeDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=url_shortener port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db

	log.Println("Database connected successfully")
}
