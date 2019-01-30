package common

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type CheckEmpty interface {
	IsEmpty() bool
}

const _SizeOfAddress = unsafe.Sizeof(uintptr(0))

func Empty(i interface{}) bool {

	if *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&i)) + _SizeOfAddress)) == 0 {
		return true
	}

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
	case CheckEmpty:
		return t.IsEmpty()
	default:
		return false
	}
}

/*func Empty(i interface{}) bool {
	if i == nil {
		return true
	}
	current := reflect.ValueOf(i).Interface()
	empty := reflect.Zero(reflect.ValueOf(i).Type()).Interface()

	return reflect.DeepEqual(current, empty)
}*/

// FileExist checks if a file exists at filePath.
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// AbsolutePath returns datadir + filename, or filename if it is absolute.
func AbsolutePath(dataDir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(dataDir, filename)
}

func WriteUint8(i uint8) []byte {
	return []byte{byte(i)}
}

func WriteUint16(i uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, i)
	return buf
}

func WriteUint32(i uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i)
	return buf
}

func WriteUint64(i uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i)
	return buf
}

func WriteInt8(i int8) []byte {
	return WriteUint8(uint8(i))
}

func WriteInt16(i int16) []byte {
	return WriteUint16(uint16(i))
}

func WriteInt32(i int32) []byte {
	return WriteUint32(uint32(i))
}

func WriteInt64(i int64) []byte {
	return WriteUint64(uint64(i))
}

//
//func WriteString(s string) []byte {
//	return WriteByteArray([]byte(s))
//}
//
//func WriteByteArray(b []byte) []byte {
//	//EosAssert(len(b) <= MAX_SIZE_OF_BYTE_ARRAYS, &exception.AssertException{}, "rlp encode ByteArray")
//	if err := WriteUVarInt(len(b)); err != nil {
//		return err
//	}
//	return e.toWriter(b)
//}
func WriteUVarInt(v int) []byte {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return buf[:l]
}
func WriteVarInt(v int) []byte {
	buf := make([]byte, 8)
	l := binary.PutVarint(buf, int64(v))
	return buf[:l]
}

func ReadUvarint64(in []byte) (uint64, int, error) {
	l, read := binary.Uvarint(in)
	if read < 0 {
		return l, 0, fmt.Errorf("too short")
	}

	return l, read, nil
}

func ReadVarint64(in []byte) (int64, int, error) {
	l, read := binary.Varint(in)
	if read < 0 {
		return l, 0, fmt.Errorf("too short")
	}

	return l, read, nil
}

//func  ReadVarint32() (out int32, err error) {
//	n, err := d.ReadVarint64()
//	if err != nil {
//		return out, err
//	}
//	out = int32(n)
//	return
//}
//func  ReadUvarint32() (out uint32, err error) {
//	n, err := d.ReadUvarint64()
//	if err != nil {
//		return out, err
//	}
//	out = uint32(n)
//	return
//}
