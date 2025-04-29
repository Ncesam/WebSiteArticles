package postgres

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           int64  `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex"`
	Nickname     string `gorm:"uniqueIndex"`
	HashPassword string
}
