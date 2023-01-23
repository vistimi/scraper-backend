package dynamodb

import (
	controllerModel "scraper-backend/src/adapter/controller/model"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	Type         string    `dynamodbav:"Type"` // PK
	ID           uuid.UUID `dynamodbav:"ID"`   // SK
	Name         string    `dynamodbav:"Name"`
	CreationDate time.Time `dynamodbav:"CreationDate"`
	OriginName   string    `dynamodbav:"OriginName"` // user to create tag
}

func (t *Tag) DriverMarshal(value controllerModel.Tag) {
	t.Type = value.Type
	t.ID = value.ID
	t.Name = value.Name
	t.CreationDate = value.CreationDate
	t.OriginName = value.OriginName
}

func (t Tag) DriverUnmarshal() controllerModel.Tag {
	return controllerModel.Tag{
		Type:         t.Type,
		ID:           t.ID,
		Name:         t.Name,
		CreationDate: t.CreationDate,
		OriginName:   t.OriginName,
	}
}
