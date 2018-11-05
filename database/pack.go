package database

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/eosspark/eos-go/exception"
	"reflect"
)

const (
	MAX_NUM_ARRAY_ELEMENT   = int(1024 * 1024)
	MAX_SIZE_OF_BYTE_ARRAYS = int(20 * 1024 * 1024)
)

type Pack interface {
	EncodeRLP(io.Writer) error
}

// --------------------------------------------------------------
// Encoder implements the EOS packing, similar to FC_BUFFER
// --------------------------------------------------------------
type encoder struct {
	output   io.Writer
	Order    binary.ByteOrder
	count    int
	eosArray bool
	vuint32  bool
}

var staticVariantTag uint8
var trxIsID bool

func newEncoder(w io.Writer) *encoder {
	return &encoder{
		output: w,
		Order:  binary.BigEndian,
		count:  0,
	}
}

func Encode(w io.Writer, val interface{}) error {
	encoder := newEncoder(w)
	err := encoder.encode(val)
	if err != nil {
		return err
	}
	return nil
}

func EncodeToBytes(val interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := newEncoder(buf).encode(val); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeToReader(val interface{}) (size int, r io.Reader, err error) {
	buf := new(bytes.Buffer)
	if err := newEncoder(buf).encode(val); err != nil {
		return 0, nil, err
	}
	return buf.Len(), bytes.NewReader(buf.Bytes()), nil
}

func EncodeSize(val interface{}) (int, error) {
	buffer, err := EncodeToBytes(val)
	if err != nil {
		return 0, err
	}
	return len(buffer), nil
}

func (e *encoder) encode(v interface{}) (err error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	t := rv.Type()

	switch t.Kind() {
	case reflect.String:
		return e.writeString(rv.String())
	case reflect.Bool:
		return e.writeBool(rv.Bool())
	case reflect.Int8:
		return e.writeByte(byte(rv.Int()))
	case reflect.Int16:
		return e.writeInt16(int16(rv.Int()))
	case reflect.Int32:
		return e.writeInt32(int32(rv.Int()))
	case reflect.Int:
		return e.writeInt32(int32(rv.Int()))
	case reflect.Int64:
		return e.writeInt64(rv.Int())
	case reflect.Uint8:
		return e.writeUint8(uint8(rv.Uint()))
	case reflect.Uint16:
		return e.writeUint16(uint16(rv.Uint()))
	case reflect.Uint32:
		return e.writeUint32(uint32(rv.Uint()))
	case reflect.Uint:
		return e.writeUint32(uint32(rv.Uint()))
	case reflect.Uint64:
		return e.writeUint64(rv.Uint())

	case reflect.Array:
		l := t.Len()
		exception.EosAssert(l <= MAX_NUM_ARRAY_ELEMENT, &exception.AssertException{}, "the length of array is too big")
		if !e.eosArray {
			if err = e.writeUVarInt(l); err != nil {
				return
			}
		}
		e.eosArray = false //normal array like [4]int need length of array

		for i := 0; i < l; i++ {
			if err = e.encode(rv.Index(i).Interface()); err != nil {
				return
			}
		}
	case reflect.Slice:
		l := rv.Len()
		exception.EosAssert(l <= MAX_NUM_ARRAY_ELEMENT, &exception.AssertException{}, "the length of slice is too big")
		if err = e.writeUVarInt(l); err != nil {
			return
		}

		for i := 0; i < l; i++ {
			if err = e.encode(rv.Index(i).Interface()); err != nil {
				return
			}
		}

	default:
		return errors.New("Encode: unsupported type " + t.String())
	}

	return
}

func (e *encoder) writeByteArray(b []byte) error {
	exception.EosAssert(len(b) <= MAX_SIZE_OF_BYTE_ARRAYS, &exception.AssertException{}, "rlp encode ByteArray")
	if err := e.writeUVarInt(len(b)); err != nil {
		return err
	}
	return e.toWriter(b)
}

func (e *encoder) writeUVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}

func (e *encoder) writeBool(b bool) (err error) {
	var out byte
	if b {
		out = 1
	}
	return e.writeByte(out)
}

func (e *encoder) writeByte(b byte) (err error) {
	return e.toWriter([]byte{b})
}

func (e *encoder) writeUint8(i uint8) (err error) {
	return e.toWriter([]byte{byte(i)})
}

func (e *encoder) writeUint16(i uint16) (err error) {
	buf := make([]byte, TypeSize.UInt16)
	binary.BigEndian.PutUint16(buf, i)
	return e.toWriter(buf)
}

func (e *encoder) writeUint32(i uint32) (err error) {
	buf := make([]byte, TypeSize.UInt32)
	binary.BigEndian.PutUint32(buf, i)
	return e.toWriter(buf)
}

func (e *encoder) writeUint64(i uint64) (err error) {
	buf := make([]byte, TypeSize.UInt64)
	binary.BigEndian.PutUint64(buf, i)
	return e.toWriter(buf)
}

func (e *encoder) writeInt8(i int8) (err error) {
	return e.writeUint8(uint8(i))
}

func (e *encoder) writeInt16(i int16) (err error) {
	return e.writeUint16(uint16(i))
}

func (e *encoder) writeInt32(i int32) (err error) {
	return e.writeUint32(uint32(i))
}

func (e *encoder) writeInt64(i int64) (err error) {
	return e.writeUint64(uint64(i))
}

func (e *encoder) writeString(s string) (err error) {
	return e.writeByteArray([]byte(s))
}

func (e *encoder) toWriter(bytes []byte) (err error) {
	e.count += len(bytes)
	_, err = e.output.Write(bytes)
	return
}
