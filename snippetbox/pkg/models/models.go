package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time // This informatio is not used by SQL Statements. It just uses UTC_TIMESTAMP()
	Expires time.Time
}
