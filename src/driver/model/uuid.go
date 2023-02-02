package model

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}

func (u *UUID) Scan(src interface{}) error {
	return (*uuid.UUID)(u).Scan(src)
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).Value()
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func ParseUUID(s string) (UUID, error) {
	uu, err := uuid.Parse(s)
	return UUID(uu), err
}
