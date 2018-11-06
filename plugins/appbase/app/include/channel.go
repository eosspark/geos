package include

import (
	"reflect"
	"fmt"
)

//type Channel struct {
//	Signal
//}

type Signal struct {
	channels []interface{}
}


func (s *Signal) connect(f interface{}) {

	v :=reflect.ValueOf(f)
	//t :=reflect.TypeOf(f)
	if s.channels != nil {
		lastIndex :=len(s.channels)-1
		lastV := reflect.ValueOf(s.channels[lastIndex])
		if v.Type() != lastV.Type() {
			fmt.Println("input type is conflict with this type, this type is ",lastV.Type(),"|input type is",v.Type())
			return
		}
	}
	s.channels = append(s.channels, f)
	//fmt.Println(s.channels)
}



func (s *Signal) emit(data... interface{}) {
	for _, e := range s.channels {
		opv := reflect.ValueOf(e)
		opt := reflect.TypeOf(e)
		inc := opt.NumIn()

		inv := make([]reflect.Value, inc)
		for i:=0; i<inc; i++ {
			inv[i] = reflect.ValueOf(data[i])
		}
		opv.Call(inv)
	}
}

func GetChannel () *Signal{
	channel := new(Signal)
	return channel
}

/**
* Publish data to a channel.  This data is *copied* on publish.
* @param data - the data to publish
*/
func (s *Signal) Publish(data... interface{}) {
	s.emit(data...)
}

/**
* subscribe to data on a channel
* @tparam Callback the type of the callback (functor|lambda)
* @param cb the callback
* @return handle to the subscription
*/
func (s *Signal) Subscribe(f interface{}) {
	s.connect(f)
}
