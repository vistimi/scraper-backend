package controller

import "time"

type User struct {
	Origin       string 
	Name         string 
	OriginID     string  
	CreationDate time.Time
}
