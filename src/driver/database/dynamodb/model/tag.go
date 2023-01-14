package dynamodb

import (
	controllerModel "scraper-backend/src/adapter/controller/model"
	"time"
)

type Tag struct {
	Type         string // PK
	Name         string // SK
	CreationDate time.Time
	OriginName   string // user to create tag
}

func (t *Tag) DriverMarshal(value controllerModel.Tag) {
	t.Type = value.Type
	t.Name = value.Name
	t.CreationDate = value.CreationDate
	t.OriginName = value.OriginName
}

func (t Tag) DriverUnmarshal() controllerModel.Tag {
	return controllerModel.Tag{
		Type:         t.Type,
		Name:         t.Name,
		CreationDate: t.CreationDate,
		OriginName:   t.OriginName,
	}
}
