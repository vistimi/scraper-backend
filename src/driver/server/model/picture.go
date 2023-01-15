package controller

import (
	"database/sql"
	"time"

	controllerModel "scraper-backend/src/adapter/controller/model"

	"github.com/google/uuid"
)

type Picture struct {
	Origin       *string                    `json:"origin,omitempty"`
	ID           *uuid.UUID                 `json:"id,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	OriginID     *string                    `json:"originID,omitempty"`
	User         *User                      `json:"user,omitempty"`
	Extension    *string                    `json:"extension,omitempty"`
	Sizes        *map[uuid.UUID]PictureSize `json:"sizes,omitempty"`
	Title        *string                    `json:"title,omitempty"`
	Description  *string                    `json:"description,omitempty"`
	License      *string                    `json:"license,omitempty"`
	CreationDate *time.Time                 `json:"creationDate,omitempty"`
	Tags         *map[uuid.UUID]PictureTag  `json:"tags,omitempty"`
}

func (p *Picture) DriverMarshal(value controllerModel.Picture) error {
	if value.Origin != "" {
		p.Origin = &value.Origin
	}
	if value.ID != uuid.Nil {
		p.ID = &value.ID
	}
	if value.OriginID != "" {
		p.OriginID = &value.OriginID
	}
	// TODO:
	p.Name = &value.Name
	p.Extension = &value.Extension
	p.Title = &value.Title
	p.Description = &value.Description
	p.License = &value.License
	p.CreationDate = &value.CreationDate

	var user User
	user.DriverMarshal(value.User)
	p.User = &user

	// size, err := ConvertMap[uuid.UUID, controllerModel.PictureSize, PictureSize](value.Size)
	// if err != nil {
	// 	return err
	// }
	sizes := make(map[uuid.UUID]PictureSize, len(value.Sizes))
	for sizeID, controllerSize := range value.Sizes {
		var driverSize PictureSize
		driverSize.DriverMarshal(controllerSize)
		sizes[sizeID] = driverSize
	}
	p.Sizes = &sizes

	// tags, err := ConvertMap[uuid.UUID, controllerModel.PictureTag, PictureTag](value.Tags)
	// if err != nil {
	// 	return err
	// }
	tags := make(map[uuid.UUID]PictureTag, len(value.Tags))
	for tagID, controllerTag := range value.Tags {
		var driverTag PictureTag
		driverTag.DriverMarshal(controllerTag)
		tags[tagID] = driverTag
	}
	p.Tags = &tags

	return nil
}

func (p Picture) DriverUnmarshal() (*controllerModel.Picture, error) {
	picture := controllerModel.Picture{}

	var length int
	if p.Sizes != nil {
		length = len(*p.Sizes)
	}
	sizes := make(map[uuid.UUID]controllerModel.PictureSize, length)
	if p.Sizes != nil {
		for sizeID, pictureSize := range *p.Sizes {
			sizes[sizeID] = PictureSize.DriverUnmarshal(pictureSize)
		}
		picture.Sizes = sizes
	}
	// size, err := ConvertMap[uuid.UUID, PictureSize, controllerModel.PictureSize](p.Size)
	// if err != nil {
	// 	return nil, err
	// }

	length = 0
	if p.Tags != nil {
		length = len(*p.Tags)
	}
	tags := make(map[uuid.UUID]controllerModel.PictureTag, length)
	if p.Tags != nil {
		for tagID, pictureTag := range *p.Tags {
			tags[tagID] = PictureTag.DriverUnmarshal(pictureTag)
		}
		picture.Tags = tags
	}

	// tags, err := ConvertMap[uuid.UUID, PictureTag, controllerModel.PictureTag](p.Tags)
	// if err != nil {
	// 	return nil, err
	// }

	// Origin:       p.Origin,
	// ID:           p.ID,
	// OriginID:     p.OriginID,
	// Name:         p.Name,
	// User:         User.DriverUnmarshal(p.User),
	// Extension:    p.Extension,

	// Title:        p.Title,
	// Description:  p.Description,
	// License:      p.License,
	// CreationDate: p.CreationDate,
	// Tags:         tags,
	return &picture, nil
}

type PictureSize struct {
	CreationDate time.Time `json:"creationDate,omitempty"`
	Box          Box       `json:"box,omitempty"` // absolut reference of the top left of new box based on the original sizes
}

func (ps *PictureSize) DriverMarshal(value controllerModel.PictureSize) {
	ps.CreationDate = value.CreationDate

	var box Box
	box.DriverMarshal(value.Box)
	ps.Box = box
}

func (ps PictureSize) DriverUnmarshal() controllerModel.PictureSize {
	return controllerModel.PictureSize{
		CreationDate: ps.CreationDate,
		Box:          Box.DriverUnmarshal(ps.Box),
	}
}

type Box struct {
	Tlx    int `json:"tlx,omitempty"`    // top left x coordinate (pointer because 0 is a possible value)
	Tly    int `json:"tly,omitempty"`    // top left y coordinate (pointer because 0 is a possible value)
	Width  int `json:"width,omitempty"`  // width (pointer because 0 is a possible value)
	Height int `json:"height,omitempty"` // height (pointer because 0 is a possible value)
}

func (b *Box) DriverMarshal(value controllerModel.Box) {
	b.Tlx = value.Tlx
	b.Tly = value.Tly
	b.Width = value.Width
	b.Height = value.Height
}

func (b Box) DriverUnmarshal() controllerModel.Box {
	return controllerModel.Box{
		Tlx:    b.Tlx,
		Tly:    b.Tly,
		Width:  b.Width,
		Height: b.Height,
	}
}

type PictureTag struct {
	Name           string          `json:"name,omitempty"`
	CreationDate   time.Time       `json:"creationDate,omitempty"`
	OriginName     string          `json:"originName,omitempty"`
	BoxInformation *BoxInformation `json:"boxInformation,omitempty"`
}

func (pt *PictureTag) DriverMarshal(value controllerModel.PictureTag) {
	pt.Name = value.Name
	pt.CreationDate = value.CreationDate
	pt.OriginName = value.OriginName

	if value.BoxInformation.Valid {
		var boxInformation BoxInformation
		boxInformation.DriverMarshal(value.BoxInformation.Body)
		pt.BoxInformation = utilModel.Nullable[BoxInformation]{
			Valid: true,
			Body:  boxInformation,
		}
	} else {
		pt.BoxInformation = utilModel.Nullable[BoxInformation]{
			Valid: false,
			Body:  BoxInformation{},
		}
	}
}

func (pt PictureTag) DriverUnmarshal() controllerModel.PictureTag {
	var boxInformation utilModel.Nullable[controllerModel.BoxInformation]
	if pt.BoxInformation.Valid {
		boxInformation = utilModel.Nullable[controllerModel.BoxInformation]{
			Valid: true,
			Body:  BoxInformation.DriverUnmarshal(pt.BoxInformation.Body),
		}
	} else {
		boxInformation = utilModel.Nullable[controllerModel.BoxInformation]{
			Valid: false,
			Body:  controllerModel.BoxInformation{},
		}
	}

	return controllerModel.PictureTag{
		Name:           pt.Name,
		CreationDate:   pt.CreationDate,
		OriginName:     pt.OriginName,
		BoxInformation: boxInformation,
	}
}

type BoxInformation struct {
	Model       sql.NullString  `json:"model,omitempty"`       // name of the model used for the detector
	Weights     sql.NullString  `json:"weights,omitempty"`     // weights of the model used for the detector
	ImageSizeID uuid.UUID       `json:"imageSizeID,omitempty"` // reference to the anchor point
	Box         Box             `json:"box,omitempty"`         // reference of the bounding box relative to the anchor
	Confidence  sql.NullFloat64 `json:"confidence,omitempty"`  // accuracy of the model
}

func (bi *BoxInformation) DriverMarshal(value controllerModel.BoxInformation) {
	bi.Model = value.Model
	bi.Weights = value.Weights
	bi.ImageSizeID = value.ImageSizeID

	var box Box
	box.DriverMarshal(value.Box)
	bi.Box = box

	bi.Confidence = value.Confidence
}

func (bi BoxInformation) DriverUnmarshal() controllerModel.BoxInformation {
	return controllerModel.BoxInformation{
		Model:       bi.Model,
		Weights:     bi.Weights,
		ImageSizeID: bi.ImageSizeID,
		Box:         Box.DriverUnmarshal(bi.Box),
		Confidence:  bi.Confidence,
	}
}
