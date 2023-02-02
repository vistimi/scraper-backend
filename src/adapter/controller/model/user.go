package controller

import (
	model "scraper-backend/src/driver/model"
	"time"
)

type User struct {
	Origin       string
	ID           model.UUID
	Name         string
	OriginID     string
	CreationDate time.Time
}
