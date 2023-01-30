package main

import (
	"log"
	"scraper-backend/src/adapter/controller"
	"scraper-backend/src/driver/server"
	"scraper-backend/src/util"
)

func main() {
	config, err := util.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	controllerPicture := controller.ConstructorPicture(*config)
	constrollerTag := controller.ConstructorTag(*config, controllerPicture)
	constrollerUser := controller.ConstructorUser(*config)
	controllerFlickr := controller.ConstructorFlickr(*config, controllerPicture, constrollerTag, constrollerUser)
	controllerPexels := controller.ConstructorPexels(*config, controllerPicture, constrollerTag, constrollerUser)
	controllerUnsplash := controller.ConstructorUnsplash(*config, controllerPicture, constrollerTag, constrollerUser)

	server := server.Contructor(controllerPicture, constrollerTag, constrollerUser, controllerFlickr, controllerPexels, controllerUnsplash)
	server.Router()
}
