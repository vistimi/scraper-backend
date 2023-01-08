package dynamodb

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	Name         string            
	CreationDate *time.Time        
	Origin       TagOrigin          // origin informations
}

type TagOrigin struct {
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`               // name of the origin `gui`, `username` or `detector`
	Model       string             `bson:"model,omitempty" json:"model,omitempty"`             // name of the model used for the detector
	Weights     string             `bson:"weights,omitempty" json:"weights,omitempty"`         // weights of the model used for the detector
	ImageSizeID uuid.UUID `bson:"imageSizeID,omitempty" json:"imageSizeID,omitempty"` // reference to the anchor point
	Box         Box                `bson:"box,omitempty" json:"box,omitempty"`                 // reference of the bounding box relative to the anchor
	Confidence  float32            `bson:"confidence,omitempty" json:"confidence,omitempty"`   // accuracy of the model
}