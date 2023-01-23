package controller

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	Type         string    
	ID           uuid.UUID 
	Name         string    
	CreationDate time.Time 
	OriginName   string    
}
