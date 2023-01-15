package dynamodb

import (
	controllerModel "scraper-backend/src/adapter/controller/model"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Origin       string    `dynamodbav:"Origin"` // PK original website
	ID           uuid.UUID `dynamodbav:"Origin"` // SK
	Name         string    `dynamodbav:"Origin"` // userName
	OriginID     string    `dynamodbav:"Origin"` // ID from the original website
	CreationDate time.Time `dynamodbav:"Origin"`
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
