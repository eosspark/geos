package rlp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"io"
	"math"
	"reflect"
)

type Pack interface {
	Pack() ([]byte, error)
}

var trxIsID bool

const (
	MAX_NUM_ARRAY_ELEMENT   = int(1024 * 1024)
	MAX_SIZE_OF_BYTE_ARRAYS = int(20 * 1024 * 1024)
)

// --------------------------------------------------------------
// Encoder implements the EOS packing, similar to FC_BUFFER
// --------------------------------------------------------------
type Encoder struct {
	output io.Writer
	count  int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		output: w,
		count:  0,
	}
}

func Encode(w io.Writer, val interface{}) error {
	encoder := NewEncoder(w)
	err := encoder.Encode(val)
	if err != nil {
		return err
	}
	return nil
}

func EncodeToBytes(val interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := NewEncoder(buf).Encode(val); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeSize(val interface{}) (int, error) {
	buffer, err := EncodeToBytes(val)
	if err != nil {
		return 0, err
	}
	return len(buffer), nil
}

func (e *Encoder) Encode(v interface{}) (err error) {
	switch p := v.(type) {
	case Pack:
		re, err := p.Pack()
		if err != nil {
			return err
		}
		return e.toWriter(re)
	case nil:
		return
	default:
		//fmt.Println("not serializer")
	}

	rv := reflect.Indirect(reflect.ValueOf(v))
	t := rv.Type()

	switch t.Kind() {
	case reflect.String:
		return e.writeString(rv.String())
	case reflect.Bool:
		return e.writeBool(rv.Bool())
	case reflect.Int8:
		return e.WriteByte(byte(rv.Int()))
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
	case reflect.Float32:
		return e.writeFloat32(float32(rv.Float()))
	case reflect.Float64:
		return e.writeFloat64(rv.Float())

	case reflect.Array:
		l := t.Len()
		EosAssert(l <= MAX_NUM_ARRAY_ELEMENT, &exception.AssertException{}, "the length of array is too big")

		for i := 0; i < l; i++ {
			if err = e.Encode(rv.Index(i).Interface()); err != nil {
				return
			}
		}
	case reflect.Slice:
		l := rv.Len()
		EosAssert(l <= MAX_NUM_ARRAY_ELEMENT, &exception.AssertException{}, "the length of slice is too big")
		if err = e.WriteUVarInt(l); err != nil {
			return
		}
		for i := 0; i < l; i++ {
			if err = e.Encode(rv.Index(i).Interface()); err != nil {
				return
			}
		}
	case reflect.Struct:
		l := rv.NumField()
		for i := 0; i < l; i++ {
			field := t.Field(i)
			tag := field.Tag.Get("eos")

			switch tag {
			case "-":
				continue
			case "optional":
				if rv.Field(i).IsNil() {
					e.writeBool(false)
					continue
				}
				e.writeBool(true)

			case "tag0":
				if rv.Field(i).IsNil() {
					e.writeUint8(0)
					trxIsID = true
					continue
				}
				e.writeUint8(1)
			case "tag1":
				if !trxIsID {
					continue
				}
			}

			if v := rv.Field(i); t.Field(i).Name != "_" {
				if v.CanInterface() {
					if err = e.Encode(v.Interface()); err != nil {
						return
					}
				}
			}

		}

	case reflect.Map:
		l := rv.Len()
		if err = e.WriteUVarInt(l); err != nil {
			return
		}
		for _, key := range rv.MapKeys() {
			value := rv.MapIndex(key)
			if err = e.Encode(key.Interface()); err != nil {
				return err
			}
			if err = e.Encode(value.Interface()); err != nil {
				return err
			}
		}

	default:
		return errors.New("Encode: unsupported type " + t.String())
	}

	return
}

func (e *Encoder) writeByteArray(b []byte) error {
	EosAssert(len(b) <= MAX_SIZE_OF_BYTE_ARRAYS, &exception.AssertException{}, "rlp encode ByteArray")
	if err := e.WriteUVarInt(len(b)); err != nil {
		return err
	}
	return e.toWriter(b)
}

func (e *Encoder) WriteUVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}
func (e *Encoder) WriteVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutVarint(buf, int64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) writeBool(b bool) (err error) {
	var out byte
	if b {
		out = 1
	}
	return e.WriteByte(out)
}

func (e *Encoder) WriteByte(b byte) (err error) {
	return e.toWriter([]byte{b})
}

func (e *Encoder) writeUint8(i uint8) (err error) {
	return e.toWriter([]byte{byte(i)})
}

func (e *Encoder) writeUint16(i uint16) (err error) {
	buf := make([]byte, TypeSize.UInt16)
	binary.LittleEndian.PutUint16(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeUint32(i uint32) (err error) {
	buf := make([]byte, TypeSize.UInt32)
	binary.LittleEndian.PutUint32(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeUint64(i uint64) (err error) {
	buf := make([]byte, TypeSize.UInt64)
	binary.LittleEndian.PutUint64(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeInt8(i int8) (err error) {
	return e.writeUint8(uint8(i))
}

func (e *Encoder) writeInt16(i int16) (err error) {
	return e.writeUint16(uint16(i))
}

func (e *Encoder) writeInt32(i int32) (err error) {
	return e.writeUint32(uint32(i))
}

func (e *Encoder) writeInt64(i int64) (err error) {
	return e.writeUint64(uint64(i))
}

func (e *Encoder) writeString(s string) (err error) {
	return e.writeByteArray([]byte(s))
}

func (e *Encoder) writeFloat32(f float32) (err error) {
	i := math.Float32bits(f)
	buf := make([]byte, TypeSize.UInt32)
	binary.LittleEndian.PutUint32(buf, i)
	return e.toWriter(buf)
}
func (e *Encoder) writeFloat64(f float64) (err error) {
	i := math.Float64bits(f)
	buf := make([]byte, TypeSize.UInt64)
	binary.LittleEndian.PutUint64(buf, i)

	return e.toWriter(buf)
}

func (e *Encoder) toWriter(bytes []byte) (err error) {
	e.count += len(bytes)
	_, err = e.output.Write(bytes)
	return
}
