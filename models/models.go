package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserName string `gorm:"unique"`
}

func (User) TableName() string {
	return "user" // change def gorm tablename to user
}

type Location struct {
	id        uint    `gorm:"primaryKey"`
	user_id   uint    `gorm:"index"`
	latitude  float64 `json:"latitude"`
	longitude float64 `json:"longitude"`
}
