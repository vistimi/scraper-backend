package controller

import (
	"time"

	controllerModel "scraper-backend/src/adapter/controller/model"

	"github.com/google/uuid"
)

type Tag struct {
	Type         *string    `json:",omitempty"`
	ID           *uuid.UUID `json:",omitempty"`
	Name         *string    `json:",omitempty"`
	CreationDate *time.Time `json:",omitempty"`
	OriginName   *string    `json:",omitempty"`
}

func (t *Tag) DriverMarshal(value controllerModel.Tag) {
	// TODO:
	t.Type = &value.Type
	t.ID = &value.ID
	t.Name = &value.Name
	t.CreationDate = &value.CreationDate
	t.OriginName = &value.OriginName
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
