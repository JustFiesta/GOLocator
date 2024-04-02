package models

type User struct {
    id       uint   `gorm:"primaryKey"`
    user_name string `gorm:"unique"`
}

type Location struct {
    id        uint    `gorm:"primaryKey"`
    user_id    uint    `gorm:"index"`
    latitude  float64 `json:"latitude"`
    longitude float64 `json:"longitude"`
}
