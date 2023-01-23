package controller

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Origin       string    
	ID           uuid.UUID 
	Name         string    
	OriginID     string    
	CreationDate time.Time 
}
