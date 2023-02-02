package dynamodb

import (
	controllerModel "scraper-backend/src/adapter/controller/model"
	model "scraper-backend/src/driver/model"
	"time"
)

type User struct {
	Origin       string     `dynamodbav:"Origin"`   // PK original website
	ID           model.UUID `dynamodbav:"ID"`       // SK
	Name         string     `dynamodbav:"Name"`     // userName
	OriginID     string     `dynamodbav:"OriginID"` // ID from the original website
	CreationDate time.Time  `dynamodbav:"CreationDate"`
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
