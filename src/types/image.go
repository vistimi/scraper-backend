package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// elements here prevent cyclical imports in other folders

// BodyUpdateImage is the body for the update of an image.
type BodyUpdateImageTagsPush struct {
	Origin string             `bson:"origin,omitempty" json:"origin,omitempty"`
	ID     primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Tags   []Tag              `bson:"tags,omitempty" json:"tags,omitempty"`
}

type BodyUpdateImageTagsPull struct {
	Origin string             `bson:"origin,omitempty" json:"origin,omitempty"`
	ID     primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Names  []string           `bson:"names,omitempty" json:"names,omitempty"`
}

type BodyImageCrop struct {
	ID   primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Box  Box                `bson:"box,omitempty" json:"box,omitempty"`
	File []byte             `bson:"file,omitempty" json:"file,omitempty"`
}

type BodyTransferImage struct {
	OriginID string `bson:"originID,omitempty" json:"originID,omitempty"`
	From     string `bson:"from,omitempty" json:"from,omitempty"`
	To       string `bson:"to,omitempty" json:"to,omitempty"`
}
