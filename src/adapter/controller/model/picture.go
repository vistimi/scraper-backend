package controller

import (
	"database/sql"
	utilModel "scraper-backend/src/util/model"
	"time"

	"github.com/google/uuid"
)

type Picture struct {
	Origin       string                    
	ID           uuid.UUID                 
	Name         string                    
	OriginID     string                    
	User         User                      
	Extension    string                    
	Sizes        map[uuid.UUID]PictureSize 
	Title        string                    
	Description  string                    
	License      string                    
	CreationDate time.Time                 
	Tags         map[uuid.UUID]PictureTag  
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
