package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
// elements here prevent cyclical imports in other folders

// BodyUpdateImage is the body for the update of an image.
type BodyUpdateImage struct {
	Collection string 				`bson:"collection" json:"collection"`
	ID         primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Tags       []Tag              `bson:"tags,omitempty" json:"tags,omitempty"`
}