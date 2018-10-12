package database

import "errors"
// Errors
var (
	// ErrNoID is returned when no ID field or id tag is found in the struct.
	ErrNoID = errors.New("database : missing struct tag id field")

	// ErrBadType is returned when a method receives an unexpected value type.
	ErrBadType = errors.New("database : provided data must be a struct or a pointer to struct")

	// ErrAlreadyExists is returned uses when trying to set an existing value on a field that has a unique index.
	ErrAlreadyExists = errors.New("database : already exists")

	// ErrUnknownTag is returned when an unexpected tag is specified.
	ErrUnknownTag = errors.New("database : unknown tag")

	// ErrIncompleteStructure is return when Some fields of an object are not assigned
	ErrIncompleteStructure = errors.New("database : Incomplete structure")

	// ErrSlicePtrNeeded is returned when an unexpected value is given, instead of a pointer to struct.
	ErrStructPtrNeeded = errors.New("database : provided target must be a pointer to struct")

	// ErrSlicePtrNeeded is returned when an unexpected value is given, instead of a pointer.
	ErrPtrNeeded = errors.New("database : provided target must be a pointer to a valid variable")

	// ErrStructNeeded is returned when an unexpected value is given, instead of a struct.
	ErrStructNeeded = errors.New("database : provided target must be a struct to a valid variable")

	// ErrNotFound is returned when the specified record is not saved in the bucket.
	ErrNotFound = errors.New("database not found")
)
