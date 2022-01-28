// Package model ...
package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Cat struct
type Cat struct {
	ID         uuid.UUID `bson:"_id"`
	Name       string    `bson:"name"`
	Age        int       `bson:"age"`
	Vaccinated bool      `bson:"vaccinated"`
}

// MarshalBinary convert struct to []byte
func (c Cat) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}
