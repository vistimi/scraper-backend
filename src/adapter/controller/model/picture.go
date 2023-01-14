package controller

import (
	"database/sql"
	utilModel "scraper-backend/src/util/model"
	"time"

	"github.com/google/uuid"
)

type Picture struct {
	Origin       string
	Name         string
	OriginID     string                    `json:",omitempty"`
	User         User                      `json:",omitempty"`
	Extension    string                    `json:",omitempty"`
	Sizes        map[uuid.UUID]PictureSize `json:",omitempty"`
	Title        string                    `json:",omitempty"`
	Description  string                    `json:",omitempty"`
	License      string                    `json:",omitempty"`
	CreationDate time.Time                 `json:",omitempty"`
	Tags         map[uuid.UUID]PictureTag  `json:",omitempty"`
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
	BoxInformation utilModel.Nullable[BoxInformation]
}

type BoxInformation struct {
	Model       sql.NullString
	Weights     sql.NullString
	ImageSizeID uuid.UUID
	Box         Box
	Confidence  sql.NullFloat64
}
