package test

import (
	"fmt"
	"strings"
	"testing"
)

/*
创建struct有两种方式
---使用new创建一个Student对象,结果为指针类型 分配内存
---使用{...}创建struct 该方式有分三种形式   分配内存
    //使用&T{...}创建struct，结果为指针类型
	//使用T{...}创建struct，结果为value类型
	//其它方式&T{key value,...} 结果为指针类型 一般不用

// := 是声明并赋值，并且系统自动推断类型，不需要var关键字 本身跟创建无关

结构的声明 没有分配内存
*/

/*
 结构的声明 没有分配内存
 */

var Student1 Student        /* 声明 Student1 为 Books 类型 */

var Student2 *Student


func StructTest01Base() {
	//structTest0101()
	//structTest0102()
	structTest0103()

}

//定义一个struct
type Student struct {
	id      int
	name    string
	address string
	age     int
}





func structTest0101() {
	//使用new创建一个Student对象,结果为指针类型
	var s *Student = new(Student)
	// := 是声明并赋值，并且系统自动推断类型，不需要var关键字 效果和
	d := Student{} //d :=new(Student) 在函数外定义会报非法错误
	d.name = "fs"
	e :=new(Student)
	e.name = "fs"
	s.id = 101
	s.name = "Mikle"
	s.address = "红旗南路"
	s.age = 18

	fmt.Printf("id:%d\n", s.id)
	fmt.Printf("name:%s\n", s.name)
	fmt.Printf("address:%s\n", s.address)
	fmt.Printf("age:%d\n", s.age)
	fmt.Println(s)
}

//创建Student的其它方式
func structTest0102() {
	//使用&T{...}创建struct，结果为指针类型
	var s1 *Student = &Student{102, "John", "Nanjing Road", 19}
	fmt.Println(s1)
	fmt.Println("modifyStudentByPointer...")
	modifyStudentByPointer(s1)
	fmt.Println(s1)

	//使用T{...}创建struct，结果为value类型
	fmt.Println("-------------")
	var s2 Student = Student{103, "Smith", "Heping Road", 20}
	fmt.Println(s2)
	fmt.Println("modifyStudent...")
	modifyStudent(s2)
	fmt.Println(s2)
	//创建并初始化一个struct时，一般使用【上述】两种方式

	//其它方式
	var s3 *Student = &Student{id: 104, name: "Lancy"}
	fmt.Printf("s3:%d,%s,%s,%d\n", s3.id, s3.name, s3.address, s3.age)
}

//struct对象属于值类型，因此需要通过函数修改其原始值的时候必须使用指针
func modifyStudent(s Student) {
	s.name = s.name + "-modify"
}
func modifyStudentByPointer(s *Student) {
	s.name = s.name + "-modify"
}

type Person struct {
	firstName string
	lastName  string
}

//使用 *Person作为参数的函数
func upPerson(p *Person) {
	p.firstName = strings.ToUpper(p.firstName)
	p.lastName = strings.ToUpper(p.lastName)
}

//调用上述方法的三种方式
func structTest0103() {
	//1- struct as a value type:
	var p1 Person
	p1.firstName = "Will"
	p1.lastName = "Smith"
	upPerson(&p1)
	fmt.Println(p1)

	//2—struct as a pointer:
	var p2 = new(Person)
	p2.firstName = "Will"
	p2.lastName = "Smith"
	(*p2).lastName = "Smith" //this is also valid
	upPerson(p2)
	fmt.Println(p2)

	//3—struct as a literal:
	var p3 = &Person{"Will", "Smith"}
	upPerson(p3)
	fmt.Println(p3)
}


type user struct {
	id int
}

func TestStruct(t *testing.T) {
	a := &user{}
	a.id = 111
	b := user{}
	b.id = 222
	c := new(user)
	c.id = 333
	fmt.Println(a, &b, c)
}