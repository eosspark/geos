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
	"strings"
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

var (
	optional           bool
	vuint32            bool
	vint32             bool
	eosArray           bool
	trxID              bool
	destaticVariantTag uint8
	rlplog             log.Logger
)

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type Decoder struct {
	data  []byte
	pos   int
	asset bool
}

func init() {
	rlplog = log.New("rlp")
	rlplog.SetHandler(log.TerminalHandler)
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

func (d *Decoder) Decode(v interface{}) (err error) {
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

	if vuint32 { //TODO
		vuint32 = false
		var r uint64
		r, _ = d.ReadUvarint64()
		rv.SetUint(r)
		return
	} else if vint32 {
		vint32 = false
		var r int64
		r, _ = d.ReadVarint64()
		rv.SetInt(r)
		return
	}

	switch t.Kind() {
	case reflect.String:
		s, err := d.ReadString()
		if err != nil {
			return err
		}
		rv.SetString(s)
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

	case reflect.Array:
		len := t.Len()

		if !eosArray {
			var l uint64
			if l, err = d.ReadUvarint64(); err != nil {
				return
			}
			if int(l) != len {
				rlplog.Warn("the l is not equal to len of array")
			}
		}
		eosArray = false

		for i := 0; i < int(len); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		var l uint64
		if l, err = d.ReadUvarint64(); err != nil {
			return
		}
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
		case "vuint32":
			vuint32 = true
		case "vint32":
			vint32 = true
		case "array":
			eosArray = true
		//	//for types.TransactionWithID !!
		case "trxID":
			destaticVariantTag, _ = d.ReadByte()
		case "tag0":
			if destaticVariantTag != 1 {
				continue
			}
		case "tag1":
			if destaticVariantTag != 0 {
				continue
			}

		case "asset":
			d.asset = true
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
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	d.pos += read
	return l, nil
}
func (d *Decoder) ReadVarint64() (out int64, err error) {
	l, read := binary.Varint(d.data[d.pos:])
	if read <= 0 {
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
	if d.asset {
		d.asset = false
		if len(d.data) < 7 {
			err = fmt.Errorf("asset symbol required [%d] bytes, remaining [%d]", 7, d.Remaining())
			return "", ErrValueTooLarge
		}
		data := d.data[d.pos : d.pos+7]
		d.pos += 7
		out = strings.TrimRight(string(data), "\x00")
		return
	}
	data, err := d.ReadByteArray()
	out = string(data)
	return
}

func (d *Decoder) ReadByte() (out byte, err error) {
	if d.Remaining() < TypeSize.Byte {
		err = fmt.Errorf("byte required [1] byte, remaining [%d]", d.Remaining())
		return
	}

	out = d.data[d.pos]
	d.pos++
	return
}

func (d *Decoder) ReadBool() (out bool, err error) {
	if d.Remaining() < TypeSize.Bool {
		err = fmt.Errorf("rlp: bool required [%d] byte, remaining [%d]", TypeSize.Bool, d.Remaining())
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
	if d.Remaining() < TypeSize.UInt8 {
		err = fmt.Errorf("rlp: byte required [1] byte, remaining [%d]", d.Remaining())
		return
	}
	out = d.data[d.pos]
	d.pos++
	return
}
func (d *Decoder) ReadUint16() (out uint16, err error) {
	if d.Remaining() < TypeSize.UInt16 {
		err = fmt.Errorf("rlp: uint16 required [%d] bytes, remaining [%d]", TypeSize.UInt16, d.Remaining())
		return
	}

	out = binary.LittleEndian.Uint16(d.data[d.pos:])
	d.pos += TypeSize.UInt16
	return
}
func (d *Decoder) ReadUint32() (out uint32, err error) {
	if d.Remaining() < TypeSize.UInt32 {
		err = fmt.Errorf("rlp: uint32 required [%d] bytes, remaining [%d]", TypeSize.UInt32, d.Remaining())
		return
	}

	out = binary.LittleEndian.Uint32(d.data[d.pos:])
	d.pos += TypeSize.UInt32
	return
}
func (d *Decoder) ReadUint() (out uint, err error) {
	if d.Remaining() < TypeSize.UInt {
		err = fmt.Errorf("rlp: uint required [%d] bytes, remaining [%d]", TypeSize.UInt, d.Remaining())
		return
	}

	out = uint(binary.LittleEndian.Uint32(d.data[d.pos:]))
	d.pos += TypeSize.UInt
	return
}
func (d *Decoder) ReadUint64() (out uint64, err error) {
	if d.Remaining() < TypeSize.UInt64 {
		err = fmt.Errorf("rlp: uint64 required [%d] bytes, remaining [%d]", TypeSize.UInt64, d.Remaining())
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
	if d.Remaining() < TypeSize.UInt128 {
		err = fmt.Errorf("%s required [%d] bytes, remaining [%d]", typeName, TypeSize.UInt128, d.Remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.UInt128]
	//out.Lo = binary.LittleEndian.Uint64(data)
	//out.Hi = binary.LittleEndian.Uint64(data[8:])

	d.pos += TypeSize.UInt128
	return data, nil
}

func (d *Decoder) ReadFloat32() (out float32, err error) {
	if d.Remaining() < TypeSize.Float32 {
		err = fmt.Errorf("float32 required [%d] bytes, remaining [%d]", TypeSize.Float32, d.Remaining())
		return
	}

	n := binary.LittleEndian.Uint32(d.data[d.pos:])
	out = math.Float32frombits(n)
	d.pos += TypeSize.Float32

	return
}

func (d *Decoder) ReadFloat64() (out float64, err error) {
	if d.Remaining() < TypeSize.Float64 {
		err = fmt.Errorf("float64 required [%d] bytes, remaining [%d]", TypeSize.Float64, d.Remaining())
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
	if d.Remaining() < TypeSize.Checksum160 {
		err = fmt.Errorf("checksum 160 required [%d] bytes, remaining [%d]", TypeSize.Checksum160, d.Remaining())
		return
	}
	out = make([]byte, TypeSize.Checksum160)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum160])
	d.pos += TypeSize.Checksum160
	return
}

func (d *Decoder) ReadChecksum256() (out []byte, err error) {
	if d.Remaining() < TypeSize.Checksum256 {
		err = fmt.Errorf("checksum 256 required [%d] bytes, remaining [%d]", TypeSize.Checksum256, d.Remaining())
		return
	}
	out = make([]byte, TypeSize.Checksum256)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum256])
	d.pos += TypeSize.Checksum256
	return
}

func (d *Decoder) ReadChecksum512() (out []byte, err error) {
	if d.Remaining() < TypeSize.Checksum512 {
		err = fmt.Errorf("checksum 512 required [%d] bytes, remaining [%d]", TypeSize.Checksum512, d.Remaining())
		return
	}
	out = make([]byte, TypeSize.Checksum512)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum512])
	d.pos += TypeSize.Checksum512
	return
}

func (d *Decoder) ReadPublicKey() (out []byte, err error) {
	if d.Remaining() < TypeSize.PublicKey {
		err = fmt.Errorf("publicKey required [%d] bytes, remaining [%d]", TypeSize.PublicKey, d.Remaining())
		return
	}
	keyContent := make([]byte, 34)
	copy(keyContent, d.data[d.pos:d.pos+TypeSize.PublicKey])

	d.pos += TypeSize.PublicKey
	return keyContent, nil
}

func (d *Decoder) ReadSignature() (out []byte, err error) {
	if d.Remaining() < TypeSize.Signature {
		err = fmt.Errorf("signature required [%d] bytes, remaining [%d]", TypeSize.Signature, d.Remaining())
		return
	}
	sigContent := make([]byte, 66)
	copy(sigContent, d.data[d.pos:d.pos+TypeSize.Signature])

	d.pos += TypeSize.Signature

	return sigContent, nil
}

//func (d *Decoder) ReadSymbol() (out *Symbol, err error) {
//
//	precision, err := d.ReadUint8()
//	if err != nil {
//		return out, fmt.Errorf("read symbol: read precision: %s", err)
//	}
//	symbol, err := d.ReadString()
//	if err != nil {
//		return out, fmt.Errorf("read symbol: read symbol: %s", err)
//	}
//
//	out = &Symbol{
//		Precision: precision,
//		Symbol:    symbol,
//	}
//	return
//}

type Symbol struct {
	Precision uint8
	Symbol    string
}

func (d *Decoder) ReadSymbol() (out *Symbol, err error) {

	precision, err := d.ReadUint8()
	if err != nil {
		return out, fmt.Errorf("read symbol: read precision: %s", err)
	}
	symbol, err := d.ReadString()
	if err != nil {
		return out, fmt.Errorf("read symbol: read symbol: %s", err)
	}

	out = &Symbol{
		Precision: precision,
		Symbol:    symbol,
	}
	return
}

type Asset struct {
	Amount int64
	Symbol
}

func (d *Decoder) ReadAsset() (out Asset, err error) {

	amount, err := d.ReadInt64()
	precision, err := d.ReadByte()
	if err != nil {
		return out, fmt.Errorf("readSymbol precision, %s", err)
	}

	if d.Remaining() < 7 {
		err = fmt.Errorf("asset symbol required [%d] bytes, remaining [%d]", 7, d.Remaining())
		return
	}

	data := d.data[d.pos : d.pos+7]
	d.pos += 7

	out = Asset{}
	out.Amount = amount
	out.Precision = precision
	out.Symbol.Symbol = strings.TrimRight(string(data), "\x00")
	return
}

type ExtendedAsset struct {
	Asset    Asset
	Contract uint64
}

func (d *Decoder) ReadExtendedAsset() (out ExtendedAsset, err error) {
	asset, err := d.ReadAsset()
	if err != nil {
		return out, fmt.Errorf("read extended asset: read asset: %s", err)
	}

	contract, err := d.ReadName()
	if err != nil {
		return out, fmt.Errorf("read extended asset: read name: %s", err)
	}

	extendedAsset := ExtendedAsset{
		Asset:    asset,
		Contract: contract,
	}

	return extendedAsset, err
}

func (d *Decoder) Remaining() int {
	return len(d.data) - d.pos
}
