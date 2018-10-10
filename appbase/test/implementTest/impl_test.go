package implementTest

import (
	"fmt"
	//"os"
	//"log"
	//"sort"
	//"gopkg.in/urfave/cli.v1"
	"sort"
	"testing"
)

type stu struct {
	name string
	age  int
	sex bool
}

func (s1 *stu) SetName(name string) {
	s1.name = name
}

func (s1 *stu) SetAge(age int) {
	s1.age = age
}

func Test_Implement(t *testing.T) {

	var stu1 = stu{"sheng", 23,false}

	var map1 map[string]human
	map1 = map[string]human{"1": &stu1}

	if key, ok := map1["1"]; ok {
		key.SetAge(3)
		fmt.Println("s" )
	} else {
		fmt.Println("不存在")
	}
	keys := make([]string, len(map1))
	i := 0
	for k, _ := range map1 {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("Key: %v, Value: %v / ", k, map1[k])
	}

}
