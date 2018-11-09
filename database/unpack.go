package database

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/common/arithmetic_types"
	"io"
	"io/ioutil"
	"reflect"
)

var (
	EOL                 = errors.New("rlp: end of list")
	ErrUnPointer        = errors.New("rlp: interface given to Decode must be a pointer")
	ErrElemTooLarge     = errors.New("rlp: element is larger than containing list")
	ErrValueTooLarge    = errors.New("rlp: value size exceeds available input length")
	ErrVarIntBufferSize = errors.New("rlp: invalid buffer size")
)

var TypeSize = struct {
	Bool        int
	Byte        int
	UInt8       int
	Int8        int
	UInt16      int
	Int16       int
	UInt32      int
	Int32       int
	UInt        int
	Int         int
	UInt64      int
	Int64       int
	SHA256Bytes int
}{
	Bool:        1,
	Byte:        1,
	UInt8:       1,
	Int8:        1,
	UInt16:      2,
	Int16:       2,
	UInt32:      4,
	Int32:       4,
	UInt:        4,
	Int:         4,
	UInt64:      8,
	Int64:       8,
	SHA256Bytes: 32,
}

var (
	optional           bool
	vuint32            bool
	eosArray           bool
	trxID              bool
	destaticVariantTag uint8
)

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type decoder struct {
	data []byte
	pos  int
}

func Decode(r io.Reader, val interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return newDecoder(data).decode(val)
}

func DecodeBytes(b []byte, val interface{}) error {
	err := newDecoder(b).decode(val)
	if err != nil {
		return err
	}
	return nil
}

func newDecoder(data []byte) *decoder {
	return &decoder{
		data: data,
		pos:  0,
	}
}

func (d *decoder) decode(v interface{}) (err error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.CanAddr() {
		return ErrUnPointer
	}
	t := rv.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		newRV := reflect.New(t)
		rv.Set(newRV)
		rv = reflect.Indirect(newRV)
	}

	if vuint32 {
		vuint32 = false
		var r uint64
		r, _ = d.readUvarint()
		rv.SetUint(r)
		return
	}

	switch v.(type) {
	case *arithmeticTypes.Float64:
		var uZ uint64
		var plus bool
		plus, err = d.readBool()
		var exp uint16
		exp, err = d.readUint16()
		var frac uint64
		frac, err = d.readUint64()

		if plus {
			exp &= 0x7FFF
			uZ = uint64(exp) << 48
			frac &= uint64(0xFFFEFFFFFFFFFFFF)
			uZ |= frac
		} else {
			uZ = uint64((uint16(0x8000)-exp)|0x8000) << 48
			uZ |= uint64(0x0001000000000000) - frac

		}
		rv.SetUint(uZ)
		return err

	case *arithmeticTypes.Float128:
		var f128 arithmeticTypes.Float128

		var plus bool
		plus, err = d.readBool()
		var exp uint16
		exp, err = d.readUint16()
		var sig64 uint64
		sig64, err = d.readUint64()
		var sig0 uint64
		sig0, err = d.readUint64()

		if plus {
			exp &= 0x7FFF
			f128.High = uint64(exp) << 48
			sig64 &= uint64(0xFFFEFFFFFFFFFFFF)
			f128.High |= sig64
			f128.Low = sig0
		} else {
			f128.High = uint64((uint16(0x8000)-exp)|0x8000) << 48
			f128.High |= uint64(0x00010000000000FE) - sig64
			f128.Low = uint64(0xFFFFFFFFFFFFFFFF) - sig0
		}
		rv.Set(reflect.ValueOf(f128))
		return err

	case *arithmeticTypes.Uint128:
		var u128 arithmeticTypes.Uint128
		u128.High, err = d.readUint64()
		u128.Low, err = d.readUint64()
		rv.Set(reflect.ValueOf(u128))
		return err
	}

	switch t.Kind() {
	case reflect.String:
		s, err := d.readString()
		if err != nil {
			return err
		}
		rv.SetString(s)
		return err
	case reflect.Bool:
		var r bool
		r, err = d.readBool()
		rv.SetBool(r)
		return
	case reflect.Int:
		var n int
		n, err = d.readInt()
		rv.SetInt(int64(n))
		return
	case reflect.Int8:
		var n int8
		n, err = d.readInt8()
		rv.SetInt(int64(n))
		return
	case reflect.Int16:
		var n int16
		n, err = d.readInt16()
		rv.SetInt(int64(n))
		return
	case reflect.Int32:
		var n int32
		n, err = d.readInt32()
		rv.SetInt(int64(n))
		return
	case reflect.Int64:
		var n int64
		n, err = d.readInt64()
		rv.SetInt(int64(n))
		return
	case reflect.Uint:
		var n uint
		n, err = d.readUint()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint8:
		var n uint8
		n, err = d.readUint8()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint16:
		var n uint16
		n, err = d.readUint16()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint32:
		var n uint32
		n, err = d.readUint32()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint64:
		var n uint64
		n, err = d.readUint64()
		rv.SetUint(n)
		return

	case reflect.Array:
		len := t.Len()

		if !eosArray {
			var l uint64
			if l, err = d.readUvarint(); err != nil {
				return
			}
			if int(l) != len {
				plog.Warn("the l is not equal to len of array")
			}
		}
		eosArray = false

		for i := 0; i < int(len); i++ {
			if err = d.decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		var l uint64
		if l, err = d.readUvarint(); err != nil {
			return
		}
		rv.Set(reflect.MakeSlice(t, int(l), int(l)))
		for i := 0; i < int(l); i++ {
			if err = d.decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}

	case reflect.Map:
		var l uint64
		if l, err = d.readUvarint(); err != nil {
			return
		}
		kt := t.Key()
		vt := t.Elem()
		rv.Set(reflect.MakeMap(t))
		for i := 0; i < int(l); i++ {
			kv := reflect.Indirect(reflect.New(kt))
			if err = d.decode(kv.Addr().Interface()); err != nil {
				return
			}
			vv := reflect.Indirect(reflect.New(vt))
			if err = d.decode(vv.Addr().Interface()); err != nil {
				return
			}
			rv.SetMapIndex(kv, vv)
		}

	case reflect.Struct:
		err = d.decodeStruct(v, t, rv)
		if err != nil {
			return
		}

	default:
		return errors.New("decode, unsupported type " + t.String())
	}

	return
}

func (d *decoder) decodeStruct(v interface{}, t reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()

	for i := 0; i < l; i++ {
		switch t.Field(i).Tag.Get("eos") {
		case "-", "SVTag":
			continue
		case "optional":
			isPresent, _ := d.readByte()
			if isPresent == 0 {
				//plog.Warn("Skipping optional OptionalProducerSchedule")
				v = nil
				continue
			}
		case "vuint32":
			vuint32 = true
		case "array":
			eosArray = true
			//	//for types.TransactionWithID !!
		case "trxID":
			destaticVariantTag, _ = d.readByte()
		case "tag0":
			if destaticVariantTag != 1 {
				continue
			}
		case "tag1":
			if destaticVariantTag != 0 {
				continue
			}
		}

		if v := rv.Field(i); v.CanSet() && t.Field(i).Name != "_" {
			iface := v.Addr().Interface()
			if err = d.decode(iface); err != nil {
				return
			}
		}
	}

	return
}

func (d *decoder) readUvarint() (uint64, error) {
	l, read := binary.Uvarint(d.data[d.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	d.pos += read
	return l, nil
}

func (d *decoder) readByteArray() (out []byte, err error) {
	l, err := d.readUvarint()
	if err != nil {
		return nil, err
	}

	if len(d.data) < d.pos+int(l) {
		return nil, ErrValueTooLarge
	}

	out = d.data[d.pos : d.pos+int(l)]
	d.pos += int(l)

	return
}

func (d *decoder) readString() (out string, err error) {
	data, err := d.readByteArray()
	out = string(data)
	return
}

func (d *decoder) readByte() (out byte, err error) {
	if d.remaining() < TypeSize.Byte {
		err = fmt.Errorf("byte required [1] byte, remaining [%d]", d.remaining())
		return
	}

	out = d.data[d.pos]
	d.pos++
	return
}

func (d *decoder) readBool() (out bool, err error) {
	if d.remaining() < TypeSize.Bool {
		err = fmt.Errorf("rlp: bool required [%d] byte, remaining [%d]", TypeSize.Bool, d.remaining())
		return
	}

	b, err := d.readByte()
	if err != nil {
		err = fmt.Errorf("readBool, %s", err)
	}
	out = b != 0
	return

}
func (d *decoder) readUint8() (out byte, err error) {
	if d.remaining() < TypeSize.UInt8 {
		err = fmt.Errorf("rlp: byte required [1] byte, remaining [%d]", d.remaining())
		return
	}
	out = d.data[d.pos]
	d.pos++
	return
}
func (d *decoder) readUint16() (out uint16, err error) {
	if d.remaining() < TypeSize.UInt16 {
		err = fmt.Errorf("rlp: uint16 required [%d] bytes, remaining [%d]", TypeSize.UInt16, d.remaining())
		return
	}

	out = binary.BigEndian.Uint16(d.data[d.pos:])
	d.pos += TypeSize.UInt16
	return
}
func (d *decoder) readUint32() (out uint32, err error) {
	if d.remaining() < TypeSize.UInt32 {
		err = fmt.Errorf("rlp: uint32 required [%d] bytes, remaining [%d]", TypeSize.UInt32, d.remaining())
		return
	}

	out = binary.BigEndian.Uint32(d.data[d.pos:])
	d.pos += TypeSize.UInt32
	return
}
func (d *decoder) readUint() (out uint, err error) {
	if d.remaining() < TypeSize.UInt {
		err = fmt.Errorf("rlp: uint required [%d] bytes, remaining [%d]", TypeSize.UInt, d.remaining())
		return
	}

	out = uint(binary.BigEndian.Uint32(d.data[d.pos:]))
	d.pos += TypeSize.UInt
	return
}
func (d *decoder) readUint64() (out uint64, err error) {
	if d.remaining() < TypeSize.UInt64 {
		err = fmt.Errorf("rlp: uint64 required [%d] bytes, remaining [%d]", TypeSize.UInt64, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.UInt64]
	out = binary.BigEndian.Uint64(data)
	d.pos += TypeSize.UInt64
	return
}

func (d *decoder) readInt8() (out int8, err error) {
	n, err := d.readUint8()
	out = int8(n)
	return
}

func (d *decoder) readInt16() (out int16, err error) {
	n, err := d.readUint16()
	out = int16(n)
	return
}
func (d *decoder) readInt32() (out int32, err error) {
	n, err := d.readUint32()
	out = int32(n)
	return
}
func (d *decoder) readInt() (out int, err error) {
	n, err := d.readUint()
	out = int(n)
	return
}
func (d *decoder) readInt64() (out int64, err error) {
	n, err := d.readUint64()
	out = int64(n)
	return
}

func (d *decoder) remaining() int {
	return len(d.data) - d.pos
}
