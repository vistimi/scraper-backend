package host

import (
	interfaceHost "scraper-backend/src/driver/interface/host"
)

func ConstructorApiFlickr() interfaceHost.DriverApiFlickr {
	return &DriverApiFlickr{}
}

func ConstructorApiUnsplash() interfaceHost.DriverApiUnsplash {
	return &DriverApiUnsplash{}
}

func ConstructorApiPexels() interfaceHost.DriverApiPexels {
	return &DriverApiPexels{}
}
