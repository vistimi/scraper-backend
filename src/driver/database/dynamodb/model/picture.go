package dynamodb

import (
	"database/sql"
	"time"

	controllerModel "scraper-backend/src/adapter/controller/model"
	model "scraper-backend/src/driver/model"
)

type Picture struct {
	Origin       string        `dynamodbav:"Origin"`   // PK original werbsite
	ID           model.UUID    `dynamodbav:"ID"`       // SK
	Name         string        `dynamodbav:"Name"`     // name <originID>_time
	OriginID     string        `dynamodbav:"OriginID"` // id from original website
	User         User          `dynamodbav:"User"`
	Extension    string        `dynamodbav:"Extension"` // type of file
	Sizes        []PictureSize `dynamodbav:"Sizes"`     // size cropping history
	Title        string        `dynamodbav:"Title"`
	Description  string        `dynamodbav:"Description"` // decription of picture
	License      string        `dynamodbav:"License"`     // type of public license
	CreationDate time.Time     `dynamodbav:"CreationDate"`
	Tags         []PictureTag  `dynamodbav:"Tags"`
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

	// size, err := ConvertMap[model.UUID, controllerModel.PictureSize, PictureSize](value.Size)
	// if err != nil {
	// 	return err
	// }
	sizes := make([]PictureSize, 0, len(value.Sizes))
	for _, controllerSize := range value.Sizes {
		var driverSize PictureSize
		driverSize.DriverMarshal(controllerSize)
		sizes = append(sizes, driverSize)
	}
	p.Sizes = sizes

	// tags, err := ConvertMap[model.UUID, controllerModel.PictureTag, PictureTag](value.Tags)
	// if err != nil {
	// 	return err
	// }
	tags := make([]PictureTag, 0, len(value.Tags))
	for _, controllerTag := range value.Tags {
		var driverTag PictureTag
		driverTag.DriverMarshal(controllerTag)
		tags = append(tags, driverTag)
	}
	p.Tags = tags
}

func (p Picture) DriverUnmarshal() *controllerModel.Picture {
	sizes := make([]controllerModel.PictureSize, 0, len(p.Sizes))
	for _, pictureSize := range p.Sizes {
		sizes = append(sizes, pictureSize.DriverUnmarshal())
	}

	tags := make([]controllerModel.PictureTag, 0, len(p.Tags))
	for _, pictureTag := range p.Tags {
		tags = append(tags, pictureTag.DriverUnmarshal())
	}

	return &controllerModel.Picture{
		Origin:       p.Origin,
		ID:           p.ID,
		OriginID:     p.OriginID,
		Name:         p.Name,
		User:         p.User.DriverUnmarshal(),
		Extension:    p.Extension,
		Sizes:        sizes,
		Title:        p.Title,
		Description:  p.Description,
		License:      p.License,
		CreationDate: p.CreationDate,
		Tags:         tags,
	}
}

type PictureSize struct {
	ID           model.UUID
	CreationDate time.Time
	Box          Box // absolut reference of the top left of new box based on the original sizes
}

func (ps *PictureSize) DriverMarshal(value controllerModel.PictureSize) {
	ps.ID = value.ID
	ps.CreationDate = value.CreationDate

	var box Box
	box.DriverMarshal(value.Box)
	ps.Box = box
}

func (ps PictureSize) DriverUnmarshal() controllerModel.PictureSize {
	return controllerModel.PictureSize{
		ID:           ps.ID,
		CreationDate: ps.CreationDate,
		Box:          ps.Box.DriverUnmarshal(),
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
	ID             model.UUID
	Name           string
	CreationDate   time.Time
	OriginName     string
	BoxInformation model.Nullable[BoxInformation] // origin informations
}

func (pt *PictureTag) DriverMarshal(value controllerModel.PictureTag) {
	pt.ID = value.ID
	pt.Name = value.Name
	pt.CreationDate = value.CreationDate
	pt.OriginName = value.OriginName

	if value.BoxInformation.Valid {
		var boxInformation BoxInformation
		boxInformation.DriverMarshal(value.BoxInformation.Body)
		pt.BoxInformation = model.NewNullable(boxInformation)
	}
}

func (pt PictureTag) DriverUnmarshal() controllerModel.PictureTag {
	var boxInformation model.Nullable[controllerModel.BoxInformation]
	if pt.BoxInformation.Valid {
		boxInformation = model.NewNullable(pt.BoxInformation.Body.DriverUnmarshal())
	}

	return controllerModel.PictureTag{
		ID:             pt.ID,
		Name:           pt.Name,
		CreationDate:   pt.CreationDate,
		OriginName:     pt.OriginName,
		BoxInformation: boxInformation,
	}
}

type BoxInformation struct {
	Model         sql.NullString  // name of the model used for the detector
	Weights       sql.NullString  // weights of the model used for the detector
	PictureSizeID model.UUID      // reference to the anchor point
	Box           Box             // reference of the bounding box relative to the anchor
	Confidence    sql.NullFloat64 // accuracy of the model
}

func (bi *BoxInformation) DriverMarshal(value controllerModel.BoxInformation) {
	bi.Model = value.Model
	bi.Weights = value.Weights
	bi.PictureSizeID = value.PictureSizeID

	var box Box
	box.DriverMarshal(value.Box)
	bi.Box = box

	bi.Confidence = value.Confidence
}

func (bi BoxInformation) DriverUnmarshal() controllerModel.BoxInformation {
	return controllerModel.BoxInformation{
		Model:         bi.Model,
		Weights:       bi.Weights,
		PictureSizeID: bi.PictureSizeID,
		Box:           bi.Box.DriverUnmarshal(),
		Confidence:    bi.Confidence,
	}
}
