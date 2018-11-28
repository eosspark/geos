package include

type Signal struct {
	callers []Caller
}

type Caller interface {
	Call(data ...interface{})
}

func (s *Signal) Connect(f Caller) {
	s.callers = append(s.callers, f)
	//fmt.Println(s.channels)
}


func (s *Signal) Emit(data ...interface{}) {
	for _, f := range s.callers {
		f.Call(data...)
	}
}
