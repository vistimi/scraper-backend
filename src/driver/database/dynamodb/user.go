package dynamodb

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`           // mongodb default id
	Origin       string             `bson:"origin,omitempty" json:"origin,omitempty"`     // original website
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`         // userName
	OriginID     string             `bson:"originID,omitempty" json:"originID,omitempty"` // ID from the original website
	CreationDate *time.Time         `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
}