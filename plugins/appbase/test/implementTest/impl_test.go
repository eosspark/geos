package implementTest

import (
	"testing"
	"fmt"
	"unsafe"
)

type AbstractHuman struct {
	sex int
}

type stu struct {
	*AbstractHuman
	name string
	age  int
}

func (g *AbstractHuman) SetSex (sex int) int {
	g.sex = sex
	return sex
}



func (s1 *stu) SetSex (sex int) int {
	s1.sex = sex
	return sex
}

func (s1 *stu) SetLocal () {}

func (s1 *stu) SetName(name string) {
	s1.name = name
}

func (s1 *stu) SetAge(age int) {
	s1.age = age
}

var g = new(AbstractHuman)
var stu1 = stu{g,"sheng", 23 }

func Test_Implement(t *testing.T) {


	fmt.Println(unsafe.Sizeof(g))
	fmt.Println(unsafe.Sizeof(int(0)))

	//stu1.SetSex(10)
	//g.SetSex(2)
	//stu1.SetSex(10)
	//fmt.Println(g.sex)
	//fmt.Println(stu1.g.sex)
	//
	//
	//
	//
	//
	//var map1 map[string]human
	//map1 = map[string]human{"1": &stu1}
	//
	//if key, ok := map1["1"]; ok {
	//	key.SetAge(3)
	//	fmt.Println("s")
	//} else {
	//	fmt.Println("不存在")
	//}
	//keys := make([]string, len(map1))
	//i := 0
	//for k, _ := range map1 {
	//	keys[i] = k
	//	i++
	//}
	//sort.Strings(keys)
	//for _, k := range keys {
	//	fmt.Printf("Key: %v, Value: %v / ", k, map1[k])
	//}

}
