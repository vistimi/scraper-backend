package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

// Structure for an image strored in MongoDB
type Image struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`             // mongodb default id
	Origin         string             `bson:"origin,omitempty" json:"origin,omitempty"`       // original werbsite
	OriginID       string             `bson:"originId,omitempty" json:"originId,omitempty"`   // id from original website
	Extension      string             `bson:"extension,omitempty" json:"extension,omitempty"` // type of file
	Path           string             `bson:"path,omitempty" json:"path,omitempty"`           // relative path name <flickrId>.<extension>
	Width          int                `bson:"width,omitempty" json:"width,omitempty"`         // width of image
	Height         int                `bson:"height,omitempty" json:"height,omitempty"`       // height of image
	Title          string             `bson:"title,omitempty" json:"title,omitempty"`
	Description    string             `bson:"description,omitempty" json:"description,omitempty"`       // decription of image
	License        string             `bson:"license,omitempty" json:"license,omitempty"`               // type of public license
	CreationDate   *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
	Tags           []Tag              `bson:"tags,omitempty" json:"tags,omitempty"`
}

// Structure for a tag strored in MongoDB
type Tag struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // mongodb default id
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	Origin       string             `bson:"origin,omitempty" json:"origin,omitempty"` // origin of the tag
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
}
