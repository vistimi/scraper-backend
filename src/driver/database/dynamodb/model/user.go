package dynamodb

import (
	controllerModel "scraper-backend/src/adapter/controller/model"
	"time"
)

type User struct {
	Origin       string    `dynamodbav:"origin" json:"origin,omitempty"`     // PK original website
	Name         string    `dynamodbav:"name" json:"name,omitempty"`         // SK userName
	OriginID     string    `dynamodbav:"originID" json:"originID,omitempty"` // ID from the original website
	CreationDate time.Time `dynamodbav:"creationDate" json:"creationDate,omitempty"`
}

func (u *User) ToDriverModel(value controllerModel.User) {
	u.Origin = value.Origin
	u.Name = value.Name
	u.OriginID = value.OriginID
	u.CreationDate = value.CreationDate
}

func (u User) FromDriverModel() controllerModel.User {
	return controllerModel.User{
		Origin:       u.Origin,
		Name:         u.Name,
		OriginID:     u.OriginID,
		CreationDate: u.CreationDate,
	}
}
