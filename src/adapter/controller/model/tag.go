package controller

import "time"

type Tag struct {
	Type         string
	Name         string
	CreationDate time.Time `json:",omitempty"`
	OriginName   string    `json:",omitempty"`
}
