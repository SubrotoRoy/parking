package model

import (
	"errors"
	"time"
)

//DbConfig is used for connection to database
type DbConfig struct {
	DbUser     string
	DbPassword string
	DbName     string
	Port       string
	Host       string
}

//Person struct will be used for mapping to database table person
type Person struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	FirstName   string    `gorm:"column:first_name"`
	LastName    string    `gorm:"column:last_name"`
	PhoneNumber int       `gorm:"column:phone_number"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

//Car struct will be used for mapping to database table car
type Car struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	CarModel  string    `gorm:"column:car_model"`
	CarNumber string    `gorm:"column:car_number"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

//Slot struct will be used for mapping to databaase table slot
type Slot struct {
	ID         int  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	IsOccupied bool `gorm:"column:is_occupied"`
}

//Parking struct will be used for mapping to databaase table parking
type Parking struct {
	ID        int        `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	PersonID  int        `json:"personId" gorm:"column:person_id"`
	CarID     int        `json:"carId" gorm:"column:car_id"`
	SlotID    int        `json:"slotId" gorm:"column:slot_id"`
	EntryTime *time.Time `json:"entryTime" gorm:"column:entry_time"`
	ExitTime  *time.Time `json:"exitTime" gorm:"column:exit_time"`
	HasExited bool       `json:"-" gorm:"column:has_exited"`
	Person    Person     `json:"person" gorm:"foreignKey:PersonID"`
	Car       Car        `json:"car" gorm:"foreignKey:CarID"`
	Slot      Slot       `json:"slot" gorm:"<-:false;foreignKey:SlotID"`
}

//Errors
var (
	ErrNoSlotAvailable = errors.New("No Slots Available")
)
