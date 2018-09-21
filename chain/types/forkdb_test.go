package types

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/eosspark/eos-go/log"
)

func Test_NewForkDatabase(t *testing.T) {
	forkdb, err := NewForkDatabase("./", "forkdb.dat", true)
	if err != nil {
		t.Error(&err)
	}
	log.Debug("forkdb block state:", forkdb)
	fmt.Println("Test_NewForkDatabase run seccuss")
	defer forkdb.database.Close()
}

func Test_AddBlockState(t *testing.T) {
	var blockState = BlockState{}

	forkdb, err := NewForkDatabase("./", "forkdb.dat", true)
	if err != nil {
		t.Error(err)
	}

	b := forkdb.AddBlockState(blockState)
	/*if er != nil {
		t.Error(er)
	}*/
	log.Debug("AddBlockState return info:", b)
}

//**********************************reflect test**************************************

func Test_Exec(t *testing.T) {
	p := map[string]interface{}{}
	p["a"] = pp
	p["b"] = ppp
	getValue := reflect.ValueOf(p["b"])
	getValue.Call(nil)

	te := reflect.ValueOf(p["a"])
	params := []reflect.Value{
		reflect.ValueOf("param"),
		reflect.ValueOf(10),
	}
	te.Call(params)
}

func ppp() {
	fmt.Println("hello")
}

func pp(s string, in int) {
	in++
	fmt.Println("hello", s, in)
}

type Demo struct {
	name string
	age  int
}

func Test_o(t *testing.T) {
	d := &Demo{"hehe", 7}
	d1 := &Demo{"hehe", 9}
	a := *d
	b := *d1
	if a == b {
		fmt.Println("aaa")
	}
	fmt.Println("bbb")
}

func (d *Demo) Operation(n string, a int) *Demo {
	d.name = "test"
	d.age = a + 1
	return d
}

func (d Demo) Operation2(n string, a int) Demo {
	d.name = "test2"
	d.age = a + 2
	return d
}

func Test_op(t *testing.T) {
	d := &Demo{"hehe", 7}
	fmt.Println("原始内容：", d)
	tmp1 := d.Operation(d.name, d.age)
	fmt.Println("第一次修改，被传递内容d：", d, "返回内容被修改")
	fmt.Println("第一次修改返回值tmp1：", tmp1, "返回内容第一次被修改")
	tmp2 := d.Operation2(d.name, d.age)
	fmt.Println("第二次修改,被传递内容d：", d, "返回值依然是第一次修改内容")
	fmt.Println("第二次修改返回值：", tmp2, "新内容出现")

}

func Test_Op(t *testing.T) {
	d := &Demo{"EOS", 30}
	t1 := reflect.TypeOf(d)
	fmt.Println(t1.Name())
	t2 := reflect.ValueOf(d)

	for i := 0; i < t1.NumMethod(); i++ {
		params := []reflect.Value{
			reflect.ValueOf("param"),
			reflect.ValueOf(10),
		}
		var s = t1.Method(i).Name
		t2.MethodByName(s).Call(params)
	}
}

func Test_p(t *testing.T) {
	d := &Demo{"EOS", 30}
	t1 := reflect.ValueOf(d)
	params := []reflect.Value{
		reflect.ValueOf("param"),
		reflect.ValueOf(10),
	}
	//fmt.Println(p1.MethodByName("Operation").Call(params))
	t1.MethodByName("Operation").Call(params)

	LargeNumberNoOverflow := int64(^uint(64)>>1) / 2
	fmt.Println(LargeNumberNoOverflow)

	b := int(^uint(0)>>1) / 2
	fmt.Println(b)

	fmt.Println(int64(b) - LargeNumberNoOverflow)

}

//**********************************reflect test**************************************
