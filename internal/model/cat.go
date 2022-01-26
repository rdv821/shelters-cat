// Package model ...
package model

import "github.com/google/uuid"

// Cat struct
type Cat struct {
	ID         uuid.UUID `bson:"_id"`
	Name       string    `bson:"name"`
	Age        int       `bson:"age"`
	Vaccinated bool      `bson:"vaccinated"`
}
