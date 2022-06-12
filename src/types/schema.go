package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type Image struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FlickrId     string             `bson:"flickId" json:"flickId"`
	Path         string             `bson:"path" json:"path"`
	Width        uint               `bson:"width" json:"width"`
	Height       uint               `bson:"height" json:"height"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description" json:"description"`
	License      string             `bson:"license" json:"license"`
	Tags         []Tag
	CreationDate time.Time `bson:"creationDate" json:"creationDate"`
}

type Tag struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string             `bson:"tagName" json:"tagName"`
	Origin       string             `bson:"origin" json:"origin"`
	CreationDate time.Time          `bson:"creationDate" json:"creationDate"`
}
