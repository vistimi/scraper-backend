package util

import (
	"database/sql/driver"
	"fmt"
)

type Nullable[T any] struct {
	Valid bool
	Body T
}

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
		return fmt.Errorf("Scan: unable to scan type %T into UUID", valueType)
	}
	return nil
}

func (nt Nullable[T]) Value() (driver.Value, error) {
	return nt.Body, nil
}