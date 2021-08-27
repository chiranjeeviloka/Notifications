package main

import (
	"notification-service/internal/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// connectMySQL connects to MySQL using GORM. Replace this with code that connects to your specific database(s).
// add more functions named similar way if you have multiple databases.
// Example: connectRedis() or connectMongoDB() or connectElastic()
func connectMySQL(connString string) (*gorm.DB, error) {

	config := &gorm.Config{
		// setup GORM config.
	}

	// TIP: there is a way to silence GORM logger. This might be useful in production.
	// Once logger is silenced it will not output executed SQL statements.

	db, err := gorm.Open(mysql.Open(connString), config)
	if err != nil {
		// It is ok to fail here, because database connection is essential for this service to work!
		return nil, err
	}

	rawDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	rawDB.SetMaxIdleConns(5)
	rawDB.SetMaxOpenConns(20)
	rawDB.SetConnMaxLifetime(time.Minute * 10)

	// Migrate/Create the schema/table
	db.AutoMigrate(&model.User{})

	return db, nil
}
