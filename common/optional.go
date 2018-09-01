package common

import "reflect"

type Operator interface {
	Operator(o interface{})
}

type Optional struct {
	T interface{}
	_valid bool
	//_value reflect.Value
}

func NewOptional() Optional {
	return Optional{_valid:true}
}

func NewOptional2(i interface{}) Optional {
	return Optional{i, true}
}

func (o *Optional) valid() bool{
	return o._valid
}

func (o *Optional) Reset(){
	if o._valid {
		defer ref()
	}
	o._valid = false
}

func ref() (T interface{}){
	return reflect.ValueOf(T)
}

func (o *Optional) Distroy(){
	o.Reset()
}

