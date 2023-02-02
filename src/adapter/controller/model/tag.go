package controller

import (
	model "scraper-backend/src/driver/model"
	"time"
)

type Tag struct {
	Type         string
	ID           model.UUID
	Name         string
	CreationDate time.Time
	OriginName   string
}
