package model

import "errors"

// ErrRecordNotFound indicates the record was not found in repo.
var ErrRecordNotFound = errors.New("record not found")
