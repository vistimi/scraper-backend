package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type UUID uuid.UUID

// database
func (u *UUID) Scan(src interface{}) error {
	return (*uuid.UUID)(u).Scan(src)
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).Value()
}

// json
func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(u).String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	return (*uuid.UUID)(u).UnmarshalBinary(data)
}

// other
func NewUUID() UUID {
	return UUID(uuid.New())
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func ParseUUID(s string) (UUID, error) {
	uu, err := uuid.Parse(s)
	return UUID(uu), err
}
