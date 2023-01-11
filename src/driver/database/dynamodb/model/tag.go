package dynamodb

import (
	controllerModel "scraper-backend/src/adapter/controller/model"
	"time"
)

type Tag struct {
	Type         string    `dynamodbav:"type" json:"type,omitempty"` // PK
	Name         string    `dynamodbav:"name" json:"name,omitempty"` // SK
	CreationDate time.Time `dynamodbav:"creationDate" json:"creationDate,omitempty"`
	OriginName   string    `dynamodbav:"originName" json:"originName,omitempty"` // user to create tag
}

func (t *Tag) ToDriverModel(value controllerModel.Tag) {
	t.Type = value.Type
	t.Name = value.Name
	t.CreationDate = value.CreationDate
	t.OriginName = value.OriginName
}

func (t Tag) FromDriverModel() controllerModel.Tag {
	return controllerModel.Tag{
		Type:         t.Type,
		Name:         t.Name,
		CreationDate: t.CreationDate,
		OriginName:   t.OriginName,
	}
}
