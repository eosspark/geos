package rlp

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
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
	//decoderInterface = reflect.TypeOf(new(Decoder)).Elem()
	bigInt = reflect.TypeOf(big.Int{})
	big0   = big.NewInt(0)
)
var prefix = make([]string, 0)

var Debug bool

var print = func(s string) {
	if Debug {
		for _, s := range prefix {
			fmt.Print(s)
		}
		fmt.Print(s)
	}
}
var println = func(args ...interface{}) {
	if Debug {
		print(fmt.Sprintf("%s\n", args...))
	}
}

func Decode(r io.Reader, val interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return NewDecoder(data).decode(val)

}

func DecodeBytes(b []byte, val interface{}) error {
	err := NewDecoder(b).decode(val)
	if err != nil {
		return err
	}
	return nil
}

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type Decoder struct {
	data     []byte
	pos      int
	optional bool
	vuint32  bool
	hash     bool
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data:     data,
		pos:      0,
		optional: false,
	}
}

func (d *Decoder) decode(v interface{}) (err error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.CanAddr() {
		return ErrUnPointer
	}
	t := rv.Type()

	println(fmt.Sprintf("Decode type [%T]", v))

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		newRV := reflect.New(t)
		rv.Set(newRV)
		rv = reflect.Indirect(newRV)
	}

	if d.optional {
		d.optional = false
		println("optinoal")
		isPresent, e := d.readByte()
		if e != nil {
			err = fmt.Errorf("decode: OptionalProducerSchedule isPresent, %s", e)
			return
		}

		if isPresent == 0 {
			println("Skipping optional OptionalProducerSchedule")
			v = nil
		} else {
			err = d.decodeStruct(v, t, rv)
			if err != nil {
				return
			}
		}
		return
	} else if d.hash {
		d.hash = false
		fmt.Println("hash")

		//s,err :=readSHA256Bytes()
		//if err != nil{
		//	return
		//}

	}

	if d.vuint32 {
		d.vuint32 = false
		var r uint64
		r, _ = d.readUvarint()
		rv.SetUint(r)
		return
	}

	switch t.Kind() {
	case reflect.Array:
		print("Reading Array")
		len := t.Len()
		l, _ := d.readUvarint()
		if len != int(l) {
			fmt.Println("the length of array is wrong", len, l)
		}
		for i := 0; i < int(len); i++ {
			if err = d.decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		print("Reading Slice length ")
		var l uint64
		if l, err = d.readUvarint(); err != nil {
			return
		}
		println(fmt.Sprintf("Slice [%T] of length: %d", v, l))
		rv.Set(reflect.MakeSlice(t, int(l), int(l)))
		for i := 0; i < int(l); i++ {
			if err = d.decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}

	case reflect.Struct:
		err = d.decodeStruct(v, t, rv)
		if err != nil {
			return
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

	case reflect.String:
		s, e := d.readString()
		if e != nil {
			err = e
			return
		}
		rv.SetString(s)
		return
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

	default:
		return errors.New("decode, unsupported type " + t.String())
	}

	return
}

func (d *Decoder) decodeStruct(v interface{}, t reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()

	if Debug {
		prefix = append(prefix, "     ")
	}
	for i := 0; i < l; i++ {
		switch t.Field(i).Tag.Get("eos") {
		case "-":
			continue
		case "optional":
			d.optional = true
			// fmt.Println("276 walker", d.optional)
		case "vuint32":
			d.vuint32 = true
			// fmt.Println("276 walker", d.vuint32)
		case "hash":
			d.hash = true
		}

		if v := rv.Field(i); v.CanSet() && t.Field(i).Name != "_" {
			iface := v.Addr().Interface()
			println(fmt.Sprintf("Field name: %s", t.Field(i).Name))
			if err = d.decode(iface); err != nil {
				return
			}

		}
	}
	if Debug {
		prefix = prefix[:len(prefix)-1]
	}
	return
}

func (d *Decoder) readSHA256Bytes() (out Sha256, err error) {

	if d.remaining() < TypeSize.SHA256Bytes {
		err = fmt.Errorf("sha256 required [%d] bytes, remaining [%d]", TypeSize.SHA256Bytes, d.remaining())
		return
	}
	for i := range out.Hash_ {
		out.Hash_[i] = binary.LittleEndian.Uint64(d.data[i*8 : (i+1)*8])
	}
	d.pos += TypeSize.SHA256Bytes
	println(fmt.Sprintf("readSHA256Bytes [%s]", hex.EncodeToString(out.Bytes())))
	return
}

func (d *Decoder) readUvarint() (uint64, error) {
	l, read := binary.Uvarint(d.data[d.pos:])
	if read <= 0 {
		println(fmt.Sprintf("readUvarint [%d]", l))
		return l, ErrVarIntBufferSize
	}
	d.pos += read
	println(fmt.Sprintf("readUvarint [%d]", l))
	return l, nil
}

func (d *Decoder) readByteArray() (out []byte, err error) {
	l, err := d.readUvarint()
	if err != nil {
		return nil, err
	}

	if len(d.data) < d.pos+int(l) {
		return nil, ErrValueTooLarge
	}

	out = d.data[d.pos : d.pos+int(l)]
	d.pos += int(l)

	println(fmt.Sprintf("readByteArray [%s]", hex.EncodeToString(out)))
	return
}

func (d *Decoder) readString() (out string, err error) {
	data, err := d.readByteArray()
	out = string(data)
	println(fmt.Sprintf("readString [%s]", out))
	return
}

func (d *Decoder) readByte() (out byte, err error) {
	if d.remaining() < TypeSize.Byte {
		err = fmt.Errorf("byte required [1] byte, remaining [%d]", d.remaining())
		return
	}

	out = d.data[d.pos]
	d.pos++
	println(fmt.Sprintf("readByte [%d]", out))
	return
}

func (d *Decoder) readBool() (out bool, err error) {
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
func (d *Decoder) readUint8() (out byte, err error) {
	if d.remaining() < TypeSize.UInt8 {
		err = fmt.Errorf("rlp: byte required [1] byte, remaining [%d]", d.remaining())
		return
	}
	out = d.data[d.pos]
	d.pos++
	println(fmt.Sprintf("readUint8 [%d]", out))
	return
}
func (d *Decoder) readUint16() (out uint16, err error) {
	if d.remaining() < TypeSize.UInt16 {
		err = fmt.Errorf("rlp: uint16 required [%d] bytes, remaining [%d]", TypeSize.UInt16, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint16(d.data[d.pos:])
	d.pos += TypeSize.UInt16
	println(fmt.Sprintf("readUint16 [%d]", out))
	return
}
func (d *Decoder) readUint32() (out uint32, err error) {
	if d.remaining() < TypeSize.UInt32 {
		err = fmt.Errorf("rlp: uint32 required [%d] bytes, remaining [%d]", TypeSize.UInt32, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint32(d.data[d.pos:])
	d.pos += TypeSize.UInt32
	println(fmt.Sprintf("readUint32 [%d]", out))
	return
}
func (d *Decoder) readUint() (out uint, err error) {
	if d.remaining() < TypeSize.UInt {
		err = fmt.Errorf("rlp: uint required [%d] bytes, remaining [%d]", TypeSize.UInt, d.remaining())
		return
	}

	out = uint(binary.LittleEndian.Uint32(d.data[d.pos:]))
	d.pos += TypeSize.UInt
	println(fmt.Sprintf("readUint [%d]", out))
	return
}
func (d *Decoder) readUint64() (out uint64, err error) {
	if d.remaining() < TypeSize.UInt64 {
		err = fmt.Errorf("rlp: uint64 required [%d] bytes, remaining [%d]", TypeSize.UInt64, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.UInt64]
	out = binary.LittleEndian.Uint64(data)
	d.pos += TypeSize.UInt64
	println(fmt.Sprintf("readUint64 [%d] [%s]", out, hex.EncodeToString(data)))
	return
}

func (d *Decoder) readInt8() (out int8, err error) {
	n, err := d.readUint8()
	out = int8(n)
	return
}

func (d *Decoder) readInt16() (out int16, err error) {
	n, err := d.readUint16()
	out = int16(n)
	return
}
func (d *Decoder) readInt32() (out int32, err error) {
	n, err := d.readUint32()
	out = int32(n)
	return
}
func (d *Decoder) readInt() (out int, err error) {
	n, err := d.readUint()
	out = int(n)
	return
}
func (d *Decoder) readInt64() (out int64, err error) {
	n, err := d.readUint64()
	out = int64(n)
	return
}

func (d *Decoder) remaining() int {
	return len(d.data) - d.pos
}
