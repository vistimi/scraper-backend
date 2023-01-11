package dynamodb

import (
	"database/sql"
	"scraper-backend/src/util"
	utilModel "scraper-backend/src/util/model"
	"time"

	controllerModel "scraper-backend/src/adapter/controller/model"

	"github.com/google/uuid"
)

type Picture struct {
	Origin       string                    `dynamodbav:"origin" json:"origin,omitempty"`     // PK original werbsite
	OriginID     string                    `dynamodbav:"originID" json:"originID,omitempty"` // id from original website
	Name         string                    `dynamodbav:"name" json:"name,omitempty"`         // SK name <originID>.<extension> extension is for copies
	User         User                      `dynamodbav:"user" json:"user,omitempty"`
	Extension    string                    `dynamodbav:"extension" json:"extension,omitempty"` // type of file
	Size         map[uuid.UUID]PictureSize `dynamodbav:"size" json:"size,omitempty"`           // size cropping history
	Title        string                    `dynamodbav:"title" json:"title,omitempty"`
	Description  string                    `dynamodbav:"description" json:"description,omitempty"` // decription of picture
	License      string                    `dynamodbav:"license" json:"license,omitempty"`         // type of public license
	CreationDate time.Time                 `dynamodbav:"creationDate" json:"creationDate,omitempty"`
	Tags         map[uuid.UUID]PictureTag  `dynamodbav:"tags" json:"tags,omitempty"`
}

type PictureSize struct {
	CreationDate time.Time `dynamodbav:"creationDate" json:"creationDate,omitempty"`
	Box          Box       `dynamodbav:"box" json:"box,omitempty"` // absolut reference of the top left of new box based on the original sizes
}

type Box struct {
	Tlx    int `dynamodbav:"tlx" json:"tlx,omitempty"`       // top left x coordinate (pointer because 0 is a possible value)
	Tly    int `dynamodbav:"tly" json:"tly,omitempty"`       // top left y coordinate (pointer because 0 is a possible value)
	Width  int `dynamodbav:"width" json:"width,omitempty"`   // width (pointer because 0 is a possible value)
	Height int `dynamodbav:"height" json:"height,omitempty"` // height (pointer because 0 is a possible value)
}

type PictureTag struct {
	Name           string                             `dynamodbav:"name" json:"name,omitempty"`
	CreationDate   time.Time                          `dynamodbav:"creationDate" json:"creationDate,omitempty"`
	OriginName     string                             `dynamodbav:"originName" json:"originName,omitempty"`
	BoxInformation utilModel.Nullable[BoxInformation] `dynamodbav:"boxInformation" json:"boxInformation,omitempty"` // origin informations
}

type BoxInformation struct {
	Model       sql.NullString  `dynamodbav:"model" json:"model,omitempty"`             // name of the model used for the detector
	Weights     sql.NullString  `dynamodbav:"weights" json:"weights,omitempty"`         // weights of the model used for the detector
	ImageSizeID uuid.UUID       `dynamodbav:"imageSizeID" json:"imageSizeID,omitempty"` // reference to the anchor point
	Box         Box             `dynamodbav:"box" json:"box,omitempty"`                 // reference of the bounding box relative to the anchor
	Confidence  sql.NullFloat64 `dynamodbav:"confidence" json:"confidence,omitempty"`   // accuracy of the model
}

func (p *Picture) ToDriverModel(value controllerModel.Picture) error {
	p.Origin = value.Origin
	p.OriginID = value.OriginID
	p.Name = value.Name
	p.Extension = value.Extension
	p.Title = value.Title
	p.Description = value.Description
	p.License = value.License
	p.CreationDate = value.CreationDate

	var user User
	user.ToDriverModel(value.User)
	p.User = user

	size, err := util.ConvertMap[uuid.UUID, controllerModel.PictureSize, PictureSize](value.Size)
	if err != nil {
		return err
	}
	p.Size = size

	tags, err := util.ConvertMap[uuid.UUID, controllerModel.PictureTag, PictureTag](value.Tags)
	if err != nil {
		return err
	}
	p.Tags = tags

	return nil
}

func (p Picture) FromDriverModel() (*controllerModel.Picture, error) {
	size, err := util.ConvertMap[uuid.UUID, PictureSize, controllerModel.PictureSize](p.Size)
	if err != nil {
		return nil, err
	}
	tags, err := util.ConvertMap[uuid.UUID, PictureTag, controllerModel.PictureTag](p.Tags)
	if err != nil {
		return nil, err
	}
	return &controllerModel.Picture{
		Origin:       p.Origin,
		OriginID:     p.OriginID,
		Name:         p.Name,
		User:         User.FromDriverModel(p.User),
		Extension:    p.Extension,
		Size:         size,
		Title:        p.Title,
		Description:  p.Description,
		License:      p.License,
		CreationDate: p.CreationDate,
		Tags:         tags,
	}, nil
}
