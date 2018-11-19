package database

import "errors"

// Errors
var (
	ErrNoID = errors.New("database : missing struct tag id field")

	ErrBadType = errors.New("database : provided data must be a struct or a pointer to struct")

	ErrAlreadyExists = errors.New("database : already exists")

	ErrIdNoSort = errors.New("database : The id field cannot use great or less")

	ErrTagInvalid = errors.New("database : Invalid tag")

	ErrIncompleteStructure = errors.New("database : Incomplete structure")

	ErrStructPtrNeeded = errors.New("database : provided target must be a pointer to struct")

	ErrPtrNeeded = errors.New("database : provided target must be a pointer to a valid variable")

	ErrNotFound = errors.New("database not found")
)
