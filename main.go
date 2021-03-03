package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/SubrotoRoy/parking/datastore"
	"github.com/SubrotoRoy/parking/kafkaservice"
	"github.com/SubrotoRoy/parking/model"
	"github.com/SubrotoRoy/parking/service"
)

func main() {

	kafkaSvc := kafkaservice.NewKafkaService()

	ctx := context.Background()

	log.Println("Reader Started")

	err := initializeGormClient()
	if err != nil {
		log.Fatalf("Couldn't connect to database. ERROR: %s", err.Error())
	}
	parkingSvc := service.NewParkingManager(datastore.Repo)

	if os.Getenv("CREATE_TABLE") == "true" {
		parkingSvc.DataRepo.CreateTables()
	}
	if os.Getenv("CREATE_SLOTS") == "true" {
		parkingSvc.DataRepo.CreateParkingSlots(createSlotSlice())
	}

	for {
		msg, err := kafkaSvc.ReadFromKafka(ctx)
		if err != nil {
			log.Println("Unable to read message from Kakfa. ERROR:", err)
			continue
		}
		if len(msg.Headers) > 0 {
			log.Println("Received Message Header,", string(msg.Headers[0].Value))
		}
		log.Println("Received Message,", string(msg.Value))
		err = parkingSvc.ManageParking(msg)
		if err == nil {
			kafkaSvc.CommitMessage(ctx, msg)
		} else {
			continue
		}
	}
}

//initializeGormClient initializes the DBRepo
func initializeGormClient() error {
	log.Println("Connecting to db")
	config := model.DbConfig{
		DbUser:     os.Getenv("DBUSER"),
		DbPassword: os.Getenv("DBPASSWORD"),
		DbName:     os.Getenv("DBNAME"),
		Port:       os.Getenv("PORT"),
		Host:       os.Getenv("Host"),
	}
	datastore.Repo = &datastore.DBRepo{}
	return datastore.Repo.DBConnect(config)
}

func createSlotSlice() []model.Slot {
	var slots []model.Slot
	slotsCount, err := strconv.Atoi(os.Getenv("SLOTS_COUNT"))
	if err != nil {
		log.Println("Invalid Slots count provided, Proceeding with count = 100")
		slotsCount = 100
	}
	for i := 0; i < slotsCount; i++ {
		slots = append(slots, model.Slot{IsOccupied: false})
	}
	return slots
}
