package server

import (
	interfaceAdapter "scraper-backend/src/adapter/interface"
	interfaceServer "scraper-backend/src/driver/interface/server"
	driverServerGin "scraper-backend/src/driver/server/gin"
)

func Contructor(
	controllerPicture interfaceAdapter.ControllerPicture,
	controllerTag interfaceAdapter.ControllerTag,
	controllerUser interfaceAdapter.ControllerUser,
	controllerFlickr interfaceAdapter.ControllerFlickr,
	controllerPexels interfaceAdapter.ControllerPexels,
	controllerUnsplash interfaceAdapter.ControllerUnsplash,
) interfaceServer.DriverServerGin {
	return &driverServerGin.DriverServerGin{
		ControllerPicture:  controllerPicture,
		ControllerTag:      controllerTag,
		ControllerUser:     controllerUser,
		ControllerFlickr:   controllerFlickr,
		ControllerPexels:   controllerPexels,
		ControllerUnsplash: controllerUnsplash,
	}
}
