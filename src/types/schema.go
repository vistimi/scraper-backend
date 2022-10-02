package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

// Structure for an image strored in MongoDB
type Image struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`           // mongodb default id
	Origin       string             `bson:"origin,omitempty" json:"origin,omitempty"`     // original werbsite
	OriginID     string             `bson:"originID,omitempty" json:"originID,omitempty"` // id from original website
	User         User               `bson:"user,omitempty" json:"user,omitempty"`
	Extension    string             `bson:"extension,omitempty" json:"extension,omitempty"` // type of file
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`           // name <originID>.<extension>
	Size         []ImageSize        `bson:"size,omitempty" json:"size,omitempty"`           // size cropping history
	Title        string             `bson:"title,omitempty" json:"title,omitempty"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"` // decription of image
	License      string             `bson:"license,omitempty" json:"license,omitempty"`         // type of public license
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
	Tags         []Tag              `bson:"tags,omitempty" json:"tags,omitempty"`
}

// Structure for a tag strored in MongoDB
type Tag struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // mongodb default id
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
	Origin       TagOrigin          `bson:"origin,omitempty" json:"origin,omitempty"` // origin informations
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`           // mongodb default id
	Origin       string             `bson:"origin,omitempty" json:"origin,omitempty"`     // original website
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`         // userName
	OriginID     string             `bson:"originID,omitempty" json:"originID,omitempty"` // ID from the original website
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
}

type ImageSize struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // mongodb default id
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
	Box          Box                `bson:"box,omitempty" json:"box,omitempty"` // absolut reference of the top left of new box based on the original sizes
}

type TagOrigin struct {
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`               // name of the origin `gui`, `username` or `detector`
	Model       string             `bson:"model,omitempty" json:"model,omitempty"`             // name of the model used for the detector
	Weights     string             `bson:"weights,omitempty" json:"weights,omitempty"`         // weights of the model used for the detector
	ImageSizeID primitive.ObjectID `bson:"imageSizeID,omitempty" json:"imageSizeID,omitempty"` // reference to the anchor point
	Box         Box                `bson:"box,omitempty" json:"box,omitempty"`                 // reference of the bounding box relative to the anchor
	Confidence  float32            `bson:"confidence,omitempty" json:"confidence,omitempty"`   // accuracy of the model
}

type Box struct {
	Tlx    *int `bson:"tlx,omitempty" json:"tlx,omitempty"`       // top left x coordinate (pointer because 0 is a possible value)
	Tly    *int `bson:"tly,omitempty" json:"tly,omitempty"`       // top left y coordinate (pointer because 0 is a possible value)
	Width  *int `bson:"width,omitempty" json:"width,omitempty"`   // width (pointer because 0 is a possible value)
	Height *int `bson:"height,omitempty" json:"height,omitempty"` // height (pointer because 0 is a possible value)
}
