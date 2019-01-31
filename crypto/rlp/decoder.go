package rlp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/log"
	"io"
	"io/ioutil"
	"math"
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
	Bool   int
	Byte   int
	UInt8  int
	Int8   int
	UInt16 int
	Int16  int
	UInt32 int
	Int32  int
	UInt   int
	Int    int
	UInt64 int
	Int64  int

	UInt128        int
	Float32        int
	Float64        int
	Checksum160    int
	Checksum256    int
	Checksum512    int
	PublicKey      int
	Signature      int
	Tstamp         int
	BlockTimestamp int
	CurrencyName   int
}{
	Bool:   1,
	Byte:   1,
	UInt8:  1,
	Int8:   1,
	UInt16: 2,
	Int16:  2,
	UInt32: 4,
	Int32:  4,
	UInt:   4,
	Int:    4,
	UInt64: 8,
	Int64:  8,

	UInt128:        16,
	Float32:        4,
	Float64:        8,
	Checksum160:    20,
	Checksum256:    32,
	Checksum512:    64,
	PublicKey:      34,
	Signature:      66,
	Tstamp:         8,
	BlockTimestamp: 4,
	CurrencyName:   7,
}

var rlplog log.Logger

type Unpack interface {
	Unpack([]byte) (int, error)
}

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type Decoder struct {
	data               []byte
	pos                int
	optional           bool
	trxID              bool
	destaticVariantTag uint8
}

func init() {
	rlplog = log.New("rlp")
	rlplog.SetHandler(log.TerminalHandler)
	rlplog.SetHandler(log.DiscardHandler())
}

func Decode(r io.Reader, val interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return NewDecoder(data).Decode(val)
}

func DecodeBytes(b []byte, val interface{}) error {
	err := NewDecoder(b).Decode(val)
	if err != nil {
		return err
	}
	return nil
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data: data,
		pos:  0,
	}
}
func (d *Decoder) GetPos() int {
	return d.pos
}

func (d *Decoder) GetData() []byte {
	return d.data
}

func (d *Decoder) Decode(v interface{}) (err error) {
	if u, ok := v.(Unpack); ok {
		pos, err := u.Unpack(d.data[d.pos:])
		if err != nil {
			fmt.Println(err)
		}
		d.pos = d.pos + pos
		return err
	}

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

	switch t.Kind() {
	case reflect.String:
		s, err := d.ReadString()
		if err != nil {
			return err
		}
		rv.SetString(s)
		rlplog.Info("decode string: %s", s)
		return err
	case reflect.Bool:
		var r bool
		r, err = d.ReadBool()
		rv.SetBool(r)
		return
	case reflect.Int:
		var n int
		n, err = d.readInt()
		rv.SetInt(int64(n))
		return
	case reflect.Int8:
		var n int8
		n, err = d.ReadInt8()
		rv.SetInt(int64(n))
		return
	case reflect.Int16:
		var n int16
		n, err = d.ReadInt16()
		rv.SetInt(int64(n))
		return
	case reflect.Int32:
		var n int32
		n, err = d.ReadInt32()
		rv.SetInt(int64(n))
		return
	case reflect.Int64:
		var n int64
		n, err = d.ReadInt64()
		rv.SetInt(int64(n))
		return
	case reflect.Uint:
		var n uint
		n, err = d.ReadUint()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint8:
		var n uint8
		n, err = d.ReadUint8()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint16:
		var n uint16
		n, err = d.ReadUint16()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint32:
		var n uint32
		n, err = d.ReadUint32()
		rv.SetUint(uint64(n))
		return
	case reflect.Uint64:
		var n uint64
		n, err = d.ReadUint64()
		rv.SetUint(n)
		return
	case reflect.Float32:
		var f float32
		f, err = d.ReadFloat32()
		rv.Set(reflect.ValueOf(f))
		return
	case reflect.Float64:
		var f float64
		f, err = d.ReadFloat64()
		rv.SetFloat(f)
		return

	case reflect.Array:
		len := t.Len()
		for i := 0; i < int(len); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		var l uint64
		if l, err = d.ReadUvarint64(); err != nil {
			fmt.Println("read varUint64: ", err)
			return
		}
		rlplog.Warn("decode slice: length is %d, type: %s", l, rv.String())

		rv.Set(reflect.MakeSlice(t, int(l), int(l)))
		for i := 0; i < int(l); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}

	case reflect.Map:
		var l uint64
		if l, err = d.ReadUvarint64(); err != nil {
			return
		}
		kt := t.Key()
		vt := t.Elem()
		rv.Set(reflect.MakeMap(t))
		for i := 0; i < int(l); i++ {
			kv := reflect.Indirect(reflect.New(kt))
			if err = d.Decode(kv.Addr().Interface()); err != nil {
				return
			}
			vv := reflect.Indirect(reflect.New(vt))
			if err = d.Decode(vv.Addr().Interface()); err != nil {
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

func (d *Decoder) decodeStruct(v interface{}, t reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()
	rlplog.Warn("decode struct:   %s, length is %d", t.String(), l)
	for i := 0; i < l; i++ {
		switch t.Field(i).Tag.Get("eos") {
		case "-", "SVTag":
			continue
		case "optional":
			isPresent, _ := d.ReadByte()
			if isPresent == 0 {
				//rlplog.Warn("Skipping optional OptionalProducerSchedule")
				v = nil
				continue
			}

		case "trxID":
			d.destaticVariantTag, _ = d.ReadByte()
		case "tag0":
			if d.destaticVariantTag != 1 {
				continue
			}
		case "tag1":
			if d.destaticVariantTag != 0 {
				continue
			}

		}

		if v := rv.Field(i); v.CanSet() && t.Field(i).Name != "_" {
			iface := v.Addr().Interface()
			if err = d.Decode(iface); err != nil {
				return
			}
		}
	}

	return
}

func (d *Decoder) ReadUvarint64() (uint64, error) {
	l, read := binary.Uvarint(d.data[d.pos:])
	if read < 0 {
		return l, ErrVarIntBufferSize
	}
	d.pos += read
	return l, nil
}
func (d *Decoder) ReadVarint64() (out int64, err error) {
	l, read := binary.Varint(d.data[d.pos:])
	if read < 0 {
		return l, ErrVarIntBufferSize
	}
	d.pos += read
	return l, nil
}
func (d *Decoder) ReadVarint32() (out int32, err error) {
	n, err := d.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int32(n)
	return
}
func (d *Decoder) ReadUvarint32() (out uint32, err error) {
	n, err := d.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint32(n)
	return
}
func (d *Decoder) ReadByteArray() (out []byte, err error) {
	l, err := d.ReadUvarint64()
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

func (d *Decoder) ReadString() (out string, err error) {
	data, err := d.ReadByteArray()
	out = string(data)
	return
}

func (d *Decoder) ReadByte() (out byte, err error) {
	if d.remaining() < TypeSize.Byte {
		err = fmt.Errorf("byte required [1] byte, remaining [%d]", d.remaining())
		return
	}

	out = d.data[d.pos]
	d.pos++
	return
}

func (d *Decoder) ReadBool() (out bool, err error) {
	if d.remaining() < TypeSize.Bool {
		err = fmt.Errorf("rlp: bool required [%d] byte, remaining [%d]", TypeSize.Bool, d.remaining())
		return
	}

	b, err := d.ReadByte()
	if err != nil {
		err = fmt.Errorf("readBool, %s", err)
	}
	out = b != 0
	return

}
func (d *Decoder) ReadUint8() (out byte, err error) {
	if d.remaining() < TypeSize.UInt8 {
		err = fmt.Errorf("rlp: byte required [1] byte, remaining [%d]", d.remaining())
		return
	}
	out = d.data[d.pos]
	d.pos++
	return
}
func (d *Decoder) ReadUint16() (out uint16, err error) {
	if d.remaining() < TypeSize.UInt16 {
		err = fmt.Errorf("rlp: uint16 required [%d] bytes, remaining [%d]", TypeSize.UInt16, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint16(d.data[d.pos:])
	d.pos += TypeSize.UInt16
	return
}
func (d *Decoder) ReadUint32() (out uint32, err error) {
	if d.remaining() < TypeSize.UInt32 {
		err = fmt.Errorf("rlp: uint32 required [%d] bytes, remaining [%d]", TypeSize.UInt32, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint32(d.data[d.pos:])
	d.pos += TypeSize.UInt32
	return
}
func (d *Decoder) ReadUint() (out uint, err error) {
	if d.remaining() < TypeSize.UInt {
		err = fmt.Errorf("rlp: uint required [%d] bytes, remaining [%d]", TypeSize.UInt, d.remaining())
		return
	}

	out = uint(binary.LittleEndian.Uint32(d.data[d.pos:]))
	d.pos += TypeSize.UInt
	return
}
func (d *Decoder) ReadUint64() (out uint64, err error) {
	if d.remaining() < TypeSize.UInt64 {
		err = fmt.Errorf("rlp: uint64 required [%d] bytes, remaining [%d]", TypeSize.UInt64, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.UInt64]
	out = binary.LittleEndian.Uint64(data)
	d.pos += TypeSize.UInt64
	return
}

func (d *Decoder) ReadInt8() (out int8, err error) {
	n, err := d.ReadUint8()
	out = int8(n)
	return
}

func (d *Decoder) ReadInt16() (out int16, err error) {
	n, err := d.ReadUint16()
	out = int16(n)
	return
}
func (d *Decoder) ReadInt32() (out int32, err error) {
	n, err := d.ReadUint32()
	out = int32(n)
	return
}
func (d *Decoder) readInt() (out int, err error) {
	n, err := d.ReadUint()
	out = int(n)
	return
}
func (d *Decoder) ReadInt64() (out int64, err error) {
	n, err := d.ReadUint64()
	out = int64(n)
	return
}

func (d *Decoder) ReadUint128(typeName string) (out []byte, err error) {
	if d.remaining() < TypeSize.UInt128 {
		err = fmt.Errorf("%s required [%d] bytes, remaining [%d]", typeName, TypeSize.UInt128, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.UInt128]
	d.pos += TypeSize.UInt128
	return data, nil
}

func (d *Decoder) ReadFloat32() (out float32, err error) {
	if d.remaining() < TypeSize.Float32 {
		err = fmt.Errorf("float32 required [%d] bytes, remaining [%d]", TypeSize.Float32, d.remaining())
		return
	}

	n := binary.LittleEndian.Uint32(d.data[d.pos:])
	out = math.Float32frombits(n)
	d.pos += TypeSize.Float32

	return
}

func (d *Decoder) ReadFloat64() (out float64, err error) {
	if d.remaining() < TypeSize.Float64 {
		err = fmt.Errorf("float64 required [%d] bytes, remaining [%d]", TypeSize.Float64, d.remaining())
		return
	}

	n := binary.LittleEndian.Uint64(d.data[d.pos:])
	out = math.Float64frombits(n)
	d.pos += TypeSize.Float64
	return
}

func (d *Decoder) ReadName() (out uint64, err error) {
	n, err := d.ReadUint64()
	return n, err
}

func (d *Decoder) ReadChecksum160() (out []byte, err error) {
	if d.remaining() < TypeSize.Checksum160 {
		err = fmt.Errorf("checksum 160 required [%d] bytes, remaining [%d]", TypeSize.Checksum160, d.remaining())
		return
	}
	out = make([]byte, TypeSize.Checksum160)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum160])
	d.pos += TypeSize.Checksum160
	return
}

func (d *Decoder) ReadChecksum256() (out []byte, err error) {
	if d.remaining() < TypeSize.Checksum256 {
		err = fmt.Errorf("checksum 256 required [%d] bytes, remaining [%d]", TypeSize.Checksum256, d.remaining())
		return
	}
	out = make([]byte, TypeSize.Checksum256)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum256])
	d.pos += TypeSize.Checksum256
	return
}

func (d *Decoder) ReadChecksum512() (out []byte, err error) {
	if d.remaining() < TypeSize.Checksum512 {
		err = fmt.Errorf("checksum 512 required [%d] bytes, remaining [%d]", TypeSize.Checksum512, d.remaining())
		return
	}
	out = make([]byte, TypeSize.Checksum512)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum512])
	d.pos += TypeSize.Checksum512
	return
}

func (d *Decoder) remaining() int {
	return len(d.data) - d.pos
}
