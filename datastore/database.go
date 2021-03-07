package datastore

import (
	"fmt"
	"log"
	"time"

	"github.com/SubrotoRoy/parking/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Repo is a variable of type Repository interface
var Repo Repository

//Repository has all the database interactions
type Repository interface {
	DBConnect(config model.DbConfig) error
	CreateTables() error
	CreateParkingRecord(parking model.Parking) error
	CreateParkingSlots(slot []model.Slot) error
	GetNearestFreeSlot() (int, error)
	UnParkCar(parking model.Parking) error
	GetParkingInfo(carNumber string) (model.Parking, error)
}

//DBRepo satisfies the interface by implementing all the methods
type DBRepo struct {
	GormDB *gorm.DB
}

//DBConnect Method to connect to Db
func (dc *DBRepo) DBConnect(config model.DbConfig) error {
	var err error
	// Format DB configs to required format connect DB
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.DbUser, config.DbPassword, config.DbName, config.Port)

	dc.GormDB, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})

	if err != nil {
		log.Printf("Unable to connect DB %v", err)
	}
	log.Printf("Postgres started at %s PORT", config.Port)
	return err
}

//CreateTables method for creating the tables
func (dc *DBRepo) CreateTables() error {
	err := dc.GormDB.Debug().AutoMigrate(&model.Person{}, &model.Car{}, &model.Parking{}, &model.Slot{})
	return err
}

//CreateParkingRecord adds records to db
func (dc *DBRepo) CreateParkingRecord(parking model.Parking) error {
	return dc.GormDB.Debug().Save(&parking).Error
}

//CreateParkingSlots creates records in slots table
func (dc *DBRepo) CreateParkingSlots(slot []model.Slot) error {
	return dc.GormDB.Debug().CreateInBatches(slot, 100).Error
}

//GetNearestFreeSlot gets the id of the nearest free slot from the slots table
func (dc *DBRepo) GetNearestFreeSlot() (int, error) {
	nearestSlot := model.Slot{}
	err := dc.GormDB.Debug().Where(`is_occupied = ?`, false).Order("id ASC").Find(&nearestSlot).Error
	if err != nil || nearestSlot.ID == 0 {
		return 0, model.ErrNoSlotAvailable
	}
	err = dc.GormDB.Debug().Exec(`Update slots set is_occupied = ? where id = ?`, true, nearestSlot.ID).Error
	if err != nil {
		return 0, err
	}
	return nearestSlot.ID, err
}

//GetParkingInfo gets the parking information based on the car number
func (dc *DBRepo) GetParkingInfo(carNumber string) (model.Parking, error) {
	var car model.Car
	var parkingInfo model.Parking
	err := dc.GormDB.Debug().Where(`car_number = ?`, carNumber).Find(&car).Error
	parkingInfo.Car = car
	err = dc.GormDB.Debug().Where(`car_id = ? and has_exited = false`, car.ID).Find(&parkingInfo).Error
	err = dc.GormDB.Debug().Where(`id=?`, parkingInfo.PersonID).Find(&(parkingInfo.Person)).Error
	return parkingInfo, err
}

//UnParkCar updates the tables to remove a parked car
func (dc *DBRepo) UnParkCar(parking model.Parking) error {
	timeNow := time.Now()
	err := dc.GormDB.Debug().Exec(`Update parkings set has_exited = ?, exit_time = ? where id = ?`, true, &timeNow, parking.ID).Error
	if err != nil {
		return err
	}
	err = dc.GormDB.Debug().Exec(`Update slots set is_occupied = ? where id = ?`, false, parking.SlotID).Error
	return nil
}
