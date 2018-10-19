package database

import "errors"

// Errors
var (
	ErrNoID = errors.New("database : missing struct tag id field")

	ErrNoName = errors.New("provided target must have a name")

	ErrBadType = errors.New("database : provided data must be a struct or a pointer to struct")

	ErrAlreadyExists = errors.New("database : already exists")

	ErrIdNoSort = errors.New("database : The id field cannot use great or less")

	ErrTagInvalid = errors.New("database : Invalid tag")

	ErrSlicePtrNeeded = errors.New("provided target must be a pointer to slice")

	ErrUnknownTag = errors.New("database : unknown tag")

	ErrIncompleteStructure = errors.New("database : Incomplete structure")

	ErrStructPtrNeeded = errors.New("database : provided target must be a pointer to struct")

	ErrPtrNeeded = errors.New("database : provided target must be a pointer to a valid variable")

	ErrStructNeeded = errors.New("database : provided target must be a struct to a valid variable")

	ErrNotFound = errors.New("database not found")
)
