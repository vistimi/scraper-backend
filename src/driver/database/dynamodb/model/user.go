package dynamodb

import (
	controllerModel "scraper-backend/src/adapter/controller/model"
	"time"
)

type User struct {
	Origin       string // PK original website
	Name         string // SK userName
	OriginID     string // ID from the original website
	CreationDate time.Time
}

func (u *User) DriverMarshal(value controllerModel.User) {
	u.Origin = value.Origin
	u.Name = value.Name
	u.OriginID = value.OriginID
	u.CreationDate = value.CreationDate
}

func (u User) DriverUnmarshal() controllerModel.User {
	return controllerModel.User{
		Origin:       u.Origin,
		Name:         u.Name,
		OriginID:     u.OriginID,
		CreationDate: u.CreationDate,
	}
}
