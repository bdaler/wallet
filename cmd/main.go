package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

const (
	StateQueued     = 0
	StateInprogress = 1
	StateImported   = 2
)

type Money float64

type Promocode struct {
	gorm.Model
	Code  string
	Price Money
}

type ImportFile struct {
	gorm.Model
	Path  string
	State uint
}

var db *gorm.DB

var err error

var (
	promocodes = []Promocode{
		{Code: "cp-12312-sad12", Price: Money(123.2)},
		{Code: "cp-12qwe-ew-sad22", Price: Money(1.2)},
	}
)

func main() {
	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5433 user=pc_user dbname=pc_db sslmode=disable password=pc_password")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	needToImport, err := getImport()
	if err != nil {
		log.Print(err)
		return
	}

	needToImport.State = uint(1)
	db.Save(&needToImport)
	log.Print(db.Error)
	//insert()
}

func getImport() (*ImportFile, error) {
	var importFile ImportFile

	result := db.Unscoped().Where("state = ?", StateQueued).First(&importFile)

	if result.Error != nil {
		return nil, result.Error
	}
	return &importFile, nil
}

func insert() {
	for index := range promocodes {
		result := db.Create(&promocodes[index])
		if result.Error != nil {
			log.Print("Err with code: ", promocodes[index].Code, " price: ", promocodes[index].Price)
		}
	}
}
