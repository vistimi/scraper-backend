package dynamodb

import (
	"database/sql"
	"time"

	controllerModel "scraper-backend/src/adapter/controller/model"
	utilModel "scraper-backend/src/util/model"

	"github.com/google/uuid"
)

type Picture struct {
	Origin       string `dynamodbav:"origin" json:"origin"` // PK original werbsite
	Name         string `dynamodbav:"name" json:"name"`     // SK name <originID>_time
	OriginID     string // id from original website
	User         User
	Extension    string                    // type of file
	Sizes        map[uuid.UUID]PictureSize // size cropping history
	Title        string
	Description  string // decription of picture
	License      string // type of public license
	CreationDate time.Time
	Tags         map[uuid.UUID]PictureTag
}

func (p *Picture) DriverMarshal(value controllerModel.Picture) error {
	p.Origin = value.Origin
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

	return nil
}

func (p Picture) DriverUnmarshal() (*controllerModel.Picture, error) {
	sizes := make(map[uuid.UUID]controllerModel.PictureSize, len(p.Sizes))
	for sizeID, pictureSize := range p.Sizes {
		sizes[sizeID] = PictureSize.DriverUnmarshal(pictureSize)
	}
	// size, err := ConvertMap[uuid.UUID, PictureSize, controllerModel.PictureSize](p.Size)
	// if err != nil {
	// 	return nil, err
	// }

	tags := make(map[uuid.UUID]controllerModel.PictureTag, len(p.Tags))
	for tagID, pictureTag := range p.Tags {
		tags[tagID] = PictureTag.DriverUnmarshal(pictureTag)
	}
	// tags, err := ConvertMap[uuid.UUID, PictureTag, controllerModel.PictureTag](p.Tags)
	// if err != nil {
	// 	return nil, err
	// }

	return &controllerModel.Picture{
		Origin:       p.Origin,
		OriginID:     p.OriginID,
		Name:         p.Name,
		User:         User.DriverUnmarshal(p.User),
		Extension:    p.Extension,
		Sizes:        sizes,
		Title:        p.Title,
		Description:  p.Description,
		License:      p.License,
		CreationDate: p.CreationDate,
		Tags:         tags,
	}, nil
}

type PictureSize struct {
	CreationDate time.Time
	Box          Box // absolut reference of the top left of new box based on the original sizes
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
	Tlx    int // top left x coordinate (pointer because 0 is a possible value)
	Tly    int // top left y coordinate (pointer because 0 is a possible value)
	Width  int // width (pointer because 0 is a possible value)
	Height int // height (pointer because 0 is a possible value)
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
	Name           string
	CreationDate   time.Time
	OriginName     string
	BoxInformation utilModel.Nullable[BoxInformation] // origin informations
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
	Model       sql.NullString  // name of the model used for the detector
	Weights     sql.NullString  // weights of the model used for the detector
	ImageSizeID uuid.UUID       // reference to the anchor point
	Box         Box             // reference of the bounding box relative to the anchor
	Confidence  sql.NullFloat64 // accuracy of the model
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
