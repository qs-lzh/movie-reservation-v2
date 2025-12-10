package model

import (
	"time"
)

type User struct {
	ID             uint     `gorm:"primaryKey"`
	Name           string   `gorm:"size:64;not null;uniqueIndex"`
	HashedPassword string   `gorm:"not null"`
	Role           UserRole `gorm:"type:varchar(16);not null"`
}

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type Movie struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:100;not null;uniqueIndex"`
	Description string `gorm:"type:text"`
}

type Showtime struct {
	ID      uint      `gorm:"primaryKey"`
	MovieID uint      `gorm:"not null;index"`
	HallID  uint      `gorm:"not null;index"`
	StartAt time.Time `gorm:"not null"`

	Movie Movie `gorm:"foreignKey:MovieID"`
	Hall  Hall  `gorm:"foreignKey:HallID"`
}

type Reservation struct {
	ID         uint `gorm:"primaryKey"`
	ShowtimeID uint `gorm:"not null;index;uniqueIndex:idx_unique_ticket"`
	SeatID     uint `gorm:"not null;index;uniqueIndex:idx_unique_ticket"`
	UserID     uint `gorm:"not null;index"`

	Showtime Showtime `gorm:"foreignKey:ShowtimeID"`
	User     User     `gorm:"foreignKey:UserID"`
}

type Hall struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64;not null;uniqueIndex"`
	SeatCount int    `gorm:"not null"`
	Rows      int    `gorm:"not null;check:rows > 0"`
	Cols      int    `gorm:"not null;check:cols > 0"`
}
