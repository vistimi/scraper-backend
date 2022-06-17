package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type Image struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FlickrId     string             `bson:"flickId,omitempty" json:"flickId,omitempty"`
	Path         string             `bson:"path,omitempty" json:"path,omitempty"`
	Width        uint               `bson:"width,omitempty" json:"width,omitempty"`
	Height       uint               `bson:"height,omitempty" json:"height,omitempty"`
	Title        string             `bson:"title,omitempty" json:"title,omitempty"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	License      string             `bson:"license,omitempty" json:"license,omitempty"`
	Tags         []Tag              `bson:"tags,omitempty" json:"tags,omitempty"`
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
}

type Tag struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Origin       string             `bson:"origin" json:"origin"`
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
}
