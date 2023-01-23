package main

import (
	"log"
	"scraper-backend/src/driver/server"
	"scraper-backend/src/util"
	"scraper-backend/src/adapter/controller"
)

func main() {
	config, err := util.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	controllerPicture := controller.ConstructorPicture(*config)
	constrollerTag := controller.ConstructorTag(*config, controllerPicture)
	constrollerUser := controller.ConstructorUser(*config)

	_ = server.Contructor(controllerPicture, constrollerTag, constrollerUser)
}
