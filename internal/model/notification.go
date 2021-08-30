package model

import "gorm.io/gorm"

type User struct {
	ID         string `gorm:"primaryKey"`
	DocumentID string `gorm:"not null"`
	Email      string `gorm:"not null"`
	Username   string `gorm:"not null"`
	Password   string `gorm:"not null"`
}

type Notification struct {
	gorm.Model
	DocumentID string `gorm:"not null"`
	UserID     string `gorm:"not null"`
	Type       int    `gorm:"not null"`
	Message    string `gorm:"not null"`
	DomainID   int    `gorm:"not null"`
	HasClicked int    `gorm:"not null"`
}
