//go:generate juv -o ./model_gen1.go $GOFILE
package example

import (
	"time"
)

type Post struct {
	ID      int       `validate:"gt=0"     json:"id"`
	Created time.Time `validate:"max=11"   json:"created"`
	Title   string    `validate:"required" json:"title"`
	Body    string    `validate:"max=50"   json:"body"`
	Draft   *bool     `validate:"required" json:"draft"`
}

type Login struct {
	ID   int    `validate:"gt=0"   json:"id"`
	Pass string `validate:"max=11" json:"pass"`
}
