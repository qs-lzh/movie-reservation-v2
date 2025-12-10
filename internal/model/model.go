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
	Seat     Seat     `gorm:"foreignKey:SeatID"`
	User     User     `gorm:"foreignKey:UserID"`
}

type Hall struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64;not null;uniqueIndex"`
	SeatCount int    `gorm:"not null"`
	Rows      int    `gorm:"not null;check:rows > 0"`
	Cols      int    `gorm:"not null;check:cols > 0"`
}

type Seat struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	HallID uint `gorm:"not null;index;uniqueIndex:idx_hall_row_col" json:"hall_id"`
	Row    int  `gorm:"not null;uniqueIndex:idx_hall_row_col;check:row>0" json:"row"`
	Col    int  `gorm:"not null;uniqueIndex:idx_hall_row_col;check:col>0" json:"col"`

	Hall Hall `gorm:"foreignKey:HallID;constraint:OnDelete:CASCADE" json:"-"`
}

type ShowtimeSeatStatus string

const (
	StatusAvailable ShowtimeSeatStatus = "available"
	StatusLocked    ShowtimeSeatStatus = "locked"
	StatusSold      ShowtimeSeatStatus = "sold"
)

type ShowtimeSeat struct {
	ID         uint               `gorm:"primaryKey" json:"id"`
	ShowtimeID uint               `gorm:"not null;index;uniqueIndex:idx_showtime_seat"`
	SeatID     uint               `gorm:"not null;index;uniqueIndex:idx_showtime_seat"`
	Status     ShowtimeSeatStatus `gorm:"type:varchar(16);not null"`

	Showtime Showtime `gorm:"foreignKey:ShowtimeID;constraint:OnDelete:CASCADE"`
	Seat     Seat     `gorm:"foreignKey:SeatID;constraint:OnDelete:CASCADE"`
}
