package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/SubrotoRoy/parking/datastore"
	"github.com/SubrotoRoy/parking/kafkaservice"
	"github.com/SubrotoRoy/parking/model"
	"github.com/segmentio/kafka-go"
)

//constants to determine the origin header
const (
	PARK   = "park"
	UNPARK = "unpark"
)

//ParkingManager will be used as receiver to associate handler functions to it
type ParkingManager struct {
	DataRepo datastore.Repository
}

//NewParkingManager return an instance of ParkingManager
func NewParkingManager(db datastore.Repository) *ParkingManager {
	return &ParkingManager{DataRepo: db}
}

//ManageParking accepts the kafka message and decides the operation to be performed on it
func (p *ParkingManager) ManageParking(msg kafka.Message) error {
	decodedMessage := model.Parking{}
	err := json.Unmarshal(msg.Value, &decodedMessage)
	if err != nil {
		log.Println("Unable to unmarshall kafka message into struct. ERROR:", err)
		return err
	}
	if string(msg.Headers[0].Value) == PARK {
		slotID, err := p.DataRepo.GetNearestFreeSlot()
		if err != nil || slotID == 0 {
			return err
		}
		timeNow := time.Now()
		decodedMessage.SlotID = slotID
		decodedMessage.Slot.ID = slotID
		decodedMessage.Slot.IsOccupied = true
		decodedMessage.EntryTime = &timeNow
		err = p.DataRepo.CreateParkingRecord(decodedMessage)
		if err != nil {
			log.Println(err)
			return err
		}
	} else if string(msg.Headers[0].Value) == UNPARK {
		kafkaSvc := kafkaservice.NewKafkaService()
		decodedMessage, err := p.DataRepo.GetParkingInfo(decodedMessage.Car.CarNumber)

		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(decodedMessage)
		err = p.DataRepo.UnParkCar(decodedMessage)
		if err != nil {
			log.Println("Error in Unparking. ERROR :", err)
			return err
		}

		c := context.Background()
		err = kafkaSvc.WriteToKafka(c, "unpark", decodedMessage)
		if err != nil {
			log.Println("Unable to post to kafka, ERROR:", err)
			return err
		}
	}
	return nil
}
