package common

import (
	"os"
	"path/filepath"
	"reflect"
)

type CheckEmpty interface {
	IsEmpty() bool
}

func empty(i interface{}) bool {
	switch t := i.(type) {
	case nil:
		return true
	case uint8:
		return t == 0
	case uint16:
		return t == 0
	case uint32:
		return t == 0
	case uint64:
		return t == 0
	case int32:
		return t == 0
	case int64:
		return t == 0
	case int:
		return t == 0
	case string:
		return t == ""
	case bool:
		return !t
	case *CheckEmpty:
		return t == nil
	case CheckEmpty:
		return t.IsEmpty()
	default:
		return false
	}
}

func Empty(i interface{}) bool {
	if i == nil {
		return true
	}
	current := reflect.ValueOf(i).Interface()
	empty := reflect.Zero(reflect.ValueOf(i).Type()).Interface()

	return reflect.DeepEqual(current, empty)
}

// FileExist checks if a file exists at filePath.
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// AbsolutePath returns datadir + filename, or filename if it is absolute.
func AbsolutePath(datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(datadir, filename)
}
