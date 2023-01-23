package controller

import (
	"database/sql"
	"time"

	controllerModel "scraper-backend/src/adapter/controller/model"
	utilModel "scraper-backend/src/util/model"

	"github.com/google/uuid"
)

type Picture struct {
	Origin       string                    `json:"origin,omitempty"`
	ID           uuid.UUID                 `json:"id,omitempty"`
	Name         string                    `json:"name,omitempty"`
	OriginID     string                    `json:"originID,omitempty"`
	User         User                      `json:"user,omitempty"`
	Extension    string                    `json:"extension,omitempty"`
	Sizes        map[uuid.UUID]PictureSize `json:"sizes,omitempty"`
	Title        string                    `json:"title,omitempty"`
	Description  string                    `json:"description,omitempty"`
	License      string                    `json:"license,omitempty"`
	CreationDate time.Time                 `json:"creationDate,omitempty"`
	Tags         map[uuid.UUID]PictureTag  `json:"tags,omitempty"`
}

func (p *Picture) DriverMarshal(value controllerModel.Picture) {
	p.Origin = value.Origin
	p.ID = value.ID
	p.OriginID = value.OriginID
	p.Name = value.Name
	p.Extension = value.Extension
	p.Title = value.Title
	p.Description = value.Description
	p.License = value.License
	p.CreationDate = value.CreationDate

	var user User
	user.DriverMarshal(value.User)
	p.User = user

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
	p.Sizes = sizes

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
	p.Tags = tags
}

func (p Picture) DriverUnmarshal() (*controllerModel.Picture) {
	picture := controllerModel.Picture{}

	sizes := make(map[uuid.UUID]controllerModel.PictureSize, len(p.Sizes))
	if p.Sizes != nil {
		for sizeID, pictureSize := range p.Sizes {
			sizes[sizeID] = pictureSize.DriverUnmarshal()
		}
		picture.Sizes = sizes
	}
	// size, err := ConvertMap[uuid.UUID, PictureSize, controllerModel.PictureSize](p.Size)
	// if err != nil {
	// 	return nil, err
	// }

	tags := make(map[uuid.UUID]controllerModel.PictureTag, len(p.Tags))
	if p.Tags != nil {
		for tagID, pictureTag := range p.Tags {
			tags[tagID] = pictureTag.DriverUnmarshal()
		}
		picture.Tags = tags
	}
	// tags, err := ConvertMap[uuid.UUID, PictureTag, controllerModel.PictureTag](p.Tags)
	// if err != nil {
	// 	return nil, err
	// }

	picture.Origin = p.Origin
	picture.ID = p.ID
	picture.OriginID = p.OriginID
	picture.Name = p.Name
	picture.User = p.User.DriverUnmarshal()
	picture.Extension = p.Extension
	picture.Title = p.Title
	picture.Description = p.Description
	picture.License = p.License
	picture.CreationDate = p.CreationDate
	picture.Tags = tags
	return &picture
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
		Box:          ps.Box.DriverUnmarshal(),
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

	// TODO: check when json.Marshaling if null == BoxInformation{}
	if value.BoxInformation.Valid {
		var boxInformation BoxInformation
		boxInformation.DriverMarshal(value.BoxInformation.Body)
		pt.BoxInformation = &boxInformation
	} else {
		pt.BoxInformation = nil
	}
}

func (pt PictureTag) DriverUnmarshal() controllerModel.PictureTag {
	var boxInformation utilModel.Nullable[controllerModel.BoxInformation]
	if pt.BoxInformation != nil {
		boxInformation = utilModel.Nullable[controllerModel.BoxInformation]{
			Valid: true,
			Body:  pt.BoxInformation.DriverUnmarshal(),
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
	Model       string    `json:"model,omitempty"`       // name of the model used for the detector
	Weights     string    `json:"weights,omitempty"`     // weights of the model used for the detector
	ImageSizeID uuid.UUID `json:"imageSizeID,omitempty"` // reference to the anchor point
	Box         Box       `json:"box,omitempty"`         // reference of the bounding box relative to the anchor
	Confidence  float64   `json:"confidence,omitempty"`  // accuracy of the model
}

func (bi *BoxInformation) DriverMarshal(value controllerModel.BoxInformation) {
	if value.Model.Valid {
		bi.Model = value.Model.String
	}
	if value.Weights.Valid {
		bi.Weights = value.Weights.String
	}
	bi.ImageSizeID = value.ImageSizeID

	var box Box
	box.DriverMarshal(value.Box)
	bi.Box = box

	if value.Confidence.Valid {
		bi.Confidence = value.Confidence.Float64
	}
}

func (bi BoxInformation) DriverUnmarshal() controllerModel.BoxInformation {
	var boxInformation controllerModel.BoxInformation
	if bi.Model != "" {
		boxInformation.Model = sql.NullString{
			String: bi.Model,
			Valid:  true,
		}
	}
	if bi.Weights != "" {
		boxInformation.Weights = sql.NullString{
			String: bi.Weights,
			Valid:  true,
		}
	}
	boxInformation.ImageSizeID = bi.ImageSizeID
	boxInformation.Box = bi.Box.DriverUnmarshal()
	if bi.Confidence != 0 {
		boxInformation.Confidence = sql.NullFloat64{
			Float64: bi.Confidence,
			Valid:   true,
		}
	}
	return boxInformation
}
