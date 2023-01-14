package controller

import "time"

type User struct {
	Origin       string
	Name         string
	OriginID     string    `json:",omitempty"`
	CreationDate time.Time `json:",omitempty"`
}
