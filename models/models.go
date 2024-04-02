package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserName string `gorm:"unique"`
	Locations []Location
}

func (User) TableName() string {
	return "user" // change def gorm tablename to user (corresponding to db)
}

type Location struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"index"` 
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (Location) TableName() string {
	return "location" // change def gorm table name to location (corresponding to db)
}
