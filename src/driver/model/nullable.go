package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	Valid bool
	Body  T
}

func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{
		Valid: true,
		Body:  value,
	}
}

// database
func (nt *Nullable[T]) Scan(value interface{}) error {
	var t T
	if value == nil {
		nt.Body, nt.Valid = t, false
		return nil
	}

	// TODO: check if useful
	// if reflect.TypeOf(value) != reflect.TypeOf(t){
	// 	return fmt.Errorf("Scan: unable to scan type %T into UUID", reflect.TypeOf(value))
	// }
	switch valueType := value.(type) {
	case T:
		nt.Valid = true
		nt.Body = value.(T)
	default:
		return fmt.Errorf("Scan: unable to scan type %T into Nullable", valueType)
	}
	return nil
}

func (nt Nullable[T]) Value() (driver.Value, error) {
	return nt.Body, nil
}

// json
func (u Nullable[T]) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(u.Body)
}

func (u *Nullable[T]) UnmarshalJSON(data []byte) error {
	var uu T
	err := json.Unmarshal(data, &uu)
	if err != nil {
		return err
	}
	*u = NewNullable(uu)
	return nil
}

// other
func (nt Nullable[T]) IsValid() bool {
	return nt.Valid
}
