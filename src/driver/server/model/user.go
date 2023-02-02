package controller

import (
	"time"

	controllerModel "scraper-backend/src/adapter/controller/model"
	model "scraper-backend/src/driver/model"
)

type User struct {
	Origin       string     `json:"origin,omitempty"`
	ID           model.UUID `json:"id,omitempty"`
	Name         string     `json:"name,omitempty"`
	OriginID     string     `json:"originID,omitempty"`
	CreationDate time.Time  `json:"creationDate,omitempty"`
}

func (u *User) DriverMarshal(value controllerModel.User) {
	u.Origin = value.Origin
	u.ID = value.ID
	u.Name = value.Name
	u.OriginID = value.OriginID
	u.CreationDate = value.CreationDate
}

func (u User) DriverUnmarshal() controllerModel.User {
	return controllerModel.User{
		Origin:       u.Origin,
		ID:           u.ID,
		Name:         u.Name,
		OriginID:     u.OriginID,
		CreationDate: u.CreationDate,
	}
}
