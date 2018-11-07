package include

import (
	"reflect"
	"fmt"
	"sort"
)

/*
The higher the number, the higher the execution level
*/

type Signal2 struct {
	Channels
}

type Channels []*Res


type Res struct {
	Level int
	Channel interface{}
}


func (c Channels) Len() int {    	 // overwrite Len()
	return len(c)
}
func (c Channels) Swap(i, j int){     // overwrite Swap()
	c[i], c[j] = c[j], c[i]
}
func (c Channels) Less(i, j int) bool {    // overwrite Less()
	return c[j].Level < c[i].Level
}

func (s *Signal2) Connect2(f interface{},level int) {

	v :=reflect.ValueOf(f)
	//t :=reflect.TypeOf(f)
	temp := &Res{level,f}

	if s.Channels != nil {
		lastIndex :=len(s.Channels)-1
		lastV := reflect.ValueOf(s.Channels[lastIndex].Channel)
		if v.Type() != lastV.Type() {
			fmt.Println("input type is conflicted  this type , this type is ",lastV.Type(),"|input type is",v.Type())
			return
		}
		s.Channels = append(s.Channels, temp)
		sort.Sort(s.Channels)
		//sort.Sort(sort.Reverse(s.Channels))
		return
	}
	s.Channels = append(s.Channels,temp)
}



func (s *Signal2) Emit2(data... interface{}) {
	for _, e := range s.Channels {
		opv := reflect.ValueOf(e.Channel)
		opt := reflect.TypeOf(e.Channel)
		inc := opt.NumIn()

		inv := make([]reflect.Value, inc)
		for i:=0; i<inc; i++ {
			inv[i] = reflect.ValueOf(data[i])
		}
		opv.Call(inv)
	}
}

func GetChannel2 () *Signal2{
	channel := new(Signal2)
	return channel
}


