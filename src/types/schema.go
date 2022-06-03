package types

import (
	"gopkg.in/mgo.v2/bson"
)

type FlickrImage struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	FlickrId    string        `bson:"flick_id" json:"flick_id"`
	Path        string        `bson:"path" json:"path"`
	Width       uint          `bson:"width" json:"width"`
	Height      uint          `bson:"height" json:"height"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	License     string        `bson:"license" json:"license"`
	Tags        []FlickTag
}

type FlickTag struct {
	ID   string `bson:"tag_id" json:"tag_id"`
	Name string `bson:"tag_name" json:"tag_name"`
}
