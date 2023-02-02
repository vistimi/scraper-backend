package controller

import (
	"database/sql"
	model "scraper-backend/src/driver/model"
	"time"
)

type Picture struct {
	Origin       string
	ID           model.UUID
	Name         string
	OriginID     string
	User         User
	Extension    string
	Sizes        map[model.UUID]PictureSize
	Title        string
	Description  string
	License      string
	CreationDate time.Time
	Tags         map[model.UUID]PictureTag
}

type PictureSize struct {
	CreationDate time.Time
	Box          Box
}

type Box struct {
	Tlx    int
	Tly    int
	Width  int
	Height int
}

type PictureTag struct {
	Name           string
	CreationDate   time.Time
	OriginName     string
	BoxInformation model.Nullable[BoxInformation]
}

type BoxInformation struct {
	Model       sql.NullString
	Weights     sql.NullString
	ImageSizeID model.UUID
	Box         Box
	Confidence  sql.NullFloat64
}
