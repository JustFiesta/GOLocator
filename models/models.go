package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"unique"`
	Locations []Location
}

func (User) TableName() string {
	return "user" // change def gorm tablename to user (corresponding to db)
}

type Location struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	DateTime  time.Time `json:"datetime"`
}

func (Location) TableName() string {
	return "location" // change def gorm table name to location (corresponding to db)
}

// BeforeCreate hook - executed before creating a new record
func (l *Location) BeforeCreate(tx *gorm.DB) (err error) {
	l.DateTime = time.Now() // set date_time to now
	return nil
}
