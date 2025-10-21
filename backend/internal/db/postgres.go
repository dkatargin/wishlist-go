package db

import (
	"log"

	"gorm.io/gorm"
)

var ORM *gorm.DB

func ConnectDB() (err error) {
	if ORM != nil {
		log.Println("ORM is already initialized")
		return nil
	}
	return nil