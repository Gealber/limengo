package models

import "errors"

var (
	DuplicateValueErr = errors.New("resource already exists")
	NotFoundErr       = errors.New("resource not found")
)
