package exec

import (
	//	"errors"
	"bytes"
	"fmt"
	"log"
	//"math"
	"os"
	"strings"

	"github.com/eosgo/common"
	"github.com/eosgo/control"
	"github.com/eosgo/cvm/wasm"
)

func To_string(name uint64) string {

	charmap := []byte(".12345abcdefghijklmnopqrstuvwxyz")
	tmp := name

	var bytes [13]byte

	for i := 0; i <= 12; i++ {
		var c byte
		if i == 0 {
			c = charmap[tmp&0x0f]
		} else {
			c = charmap[tmp&0x1f]
		}

		bytes[12-i] = c

		if i == 0 {
			tmp >>= 4
		} else {
			tmp >>= 5
		}

		//trim_right_dots( str );
	}

	str := string(bytes[:])
	//strings.Trim(str,".")

	return strings.Trim(str, ".")

}

func char_to_symbol(c byte) uint64 {
	if c >= 'a' && c <= 'z' {
		return uint64((c - 'a') + 6)
	}
	if c >= '1' && c <= '5' {
		return uint64((c - '1') + 1)
	}

	return 0
}

func N(str string) uint64 {

	var name uint64
	var i int

	for i = 0; i < len(str) && i < 12; i++ {
		// NOTE: char_to_symbol() returns char type, and without this explicit
		// expansion to uint64 type, the compilation fails at the point of usage
		// of string_to_name(), where the usage requires constant (compile time) expression.
		name |= (char_to_symbol(str[i]) & 0x1f) << uint(64-5*(i+1))
	}

	// The for-loop encoded up to 60 high bits into uint64 'name' variable,
	// if (strlen(str) > 12) then encode str[12] into the low (remaining)
	// 4 bits of 'name'
	if i == 12 {
		name |= char_to_symbol(str[12]) & 0x0F
	}

	return name
}

type Wasm_interface struct {
	context *control.Apply_context
	handles map[string]interface{}
	vm      *VM
}

func NewWasm_interface() *Wasm_interface {
	wasm_interface := Wasm_interface{handles: make(map[string]interface{})}

	wasm_interface.Register("eosio_assert", eosio_assert)
	wasm_interface.Register("action_data_size", action_data_size)
	wasm_interface.Register("read_action_data", read_action_data)
	wasm_interface.Register("current_time", current_time)
	wasm_interface.Register("require_auth2", require_auth2)
	wasm_interface.Register("memcpy", memcpy)
	wasm_interface.Register("printn", printn)
	wasm_interface.Register("prints", prints)

	return &wasm_interface
}

func (wasm_interface *Wasm_interface) Apply(code_id string, code []byte, context *control.Apply_context) {
	wasm_interface.context = context

	bf := bytes.NewBuffer([]byte(code))

	m, err := wasm.ReadModule(bf, importer)
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}

	// if *verify {
	// 	err = validate.VerifyModule(m)
	// 	if err != nil {
	// 		log.Fatalf("could not verify module: %v", err)
	// 	}
	// }

	if m.Export == nil {
		log.Fatalf("module has no export section")
	}

	vm, err := NewVM(m, wasm_interface)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}

	e, _ := m.Export.Entries["apply"]
	i := int64(e.Index)
	fidx := m.Function.Types[int(i)]
	ftype := m.Types.Entries[int(fidx)]

	wasm_interface.vm = vm

	args := make([]uint64, 3)
	args[0] = uint64(context.Receiver)
	args[1] = uint64(context.Code)
	args[2] = uint64(context.Action)

	o, err := vm.ExecCode(i, args[0], args[1], args[2])
	if err != nil {
		fmt.Printf("\n")
		log.Printf("err=%v", err)
	}
	if len(ftype.ReturnTypes) == 0 {
		fmt.Printf("\n")
	}
	fmt.Printf("%[1]v (%[1]T)\n", o)

}

func (wasm_interface *Wasm_interface) Register(name string, handler interface{}) bool {
	if _, ok := wasm_interface.handles[name]; ok {
		return false
	}

	wasm_interface.handles[name] = handler
	return true
}

func (wasm_interface *Wasm_interface) Add(handles map[string]interface{}) bool {
	for k, v := range handles {
		if _, ok := wasm_interface.handles[k]; !ok {
			wasm_interface.handles[k] = v
		}
	}
	return true
}

func (wasm_interface *Wasm_interface) GetHandles() map[string]interface{} {
	return wasm_interface.handles
}

func (wasm_interface *Wasm_interface) GetHandle(name string) interface{} {

	if _, ok := wasm_interface.handles[name]; ok {
		return wasm_interface.handles[name]
	}

	return nil
}

func importer(name string) (*wasm.Module, error) {
	f, err := os.Open(name + ".wasm")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, err := wasm.ReadModule(f, nil)
	if err != nil {
		return nil, err
	}
	// err = validate.VerifyModule(m)
	// if err != nil {
	// 	return nil, err
	// }
	return m, nil
}

func eosio_assert(wasm_interface *Wasm_interface, condition uint32, msg uint32) {

	fmt.Println("eosio_assert")
}

func action_data_size(wasm_interface *Wasm_interface) uint32 {

	fmt.Println("action_data_size")

	data := []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1} //("000000005c05a3e1") => '{"walker"}'
	return uint32(len(data))

}

func min(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func read_action_data(wasm_interface *Wasm_interface, memory uint32, buffer_size uint32) uint32 {

	fmt.Println("read_action_data")

	data := []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1} //("000000005c05a3e1") => '{"walker"}'

	//s = wasm_interface.context.act.data.size()
	s := len(data)
	if buffer_size == 0 {
		return uint32(s)
	}
	copy_size := min(buffer_size, uint32(s))
	copy(wasm_interface.vm.memory[memory:memory+copy_size], data)
	return copy_size

}

func current_time(wasm_interface *Wasm_interface) uint64 {

	fmt.Println("current_time")
	return 0
}

func require_auth2(wasm_interface *Wasm_interface, name common.AccountName, permission common.PermissionName) {

	fmt.Println("require_auth2")
}

func memcpy(wasm_interface *Wasm_interface, dest uint32, src uint32, length uint32) uint32 {

	fmt.Println("memcpy")
	copy(wasm_interface.vm.memory[dest:dest+length], wasm_interface.vm.memory[src:src+length])
	return length
}

func printn(wasm_interface *Wasm_interface, name uint64) {

	fmt.Println("printn")
	str := To_string(name)
	fmt.Println(str)

}

func prints(wasm_interface *Wasm_interface, str uint32) {

	fmt.Println("prints")

	var size uint32
	var i uint32
	for i = 0; i < 256; i++ {
		if wasm_interface.vm.memory[str+i] == 0 {
			break
		}
		size++
	}

	fmt.Println(string(wasm_interface.vm.memory[str : str+size]))
}
