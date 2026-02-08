package main

import (
    "log"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Note struct {
    ID      uint   `gorm:"primaryKey"`
    Title   string `gorm:"size:255"`
    Content string `gorm:"size:1000"`
}


var DB *gorm.DB


func InitDB() *gorm.DB {
    var err error
    DB, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }

    err = DB.AutoMigrate(&Note{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    log.Println("Database connected and migrated")
    return DB
}

// CloseDB closes the connection
func CloseDB(db *gorm.DB) {
    sqlDB, err := db.DB()
    if err != nil {
        log.Println("Failed to get SQL DB:", err)
        return
    }
    sqlDB.Close()
}
