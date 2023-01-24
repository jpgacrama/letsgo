package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Record struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires string
}
