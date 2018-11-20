package include

type Signal struct {
	functions []Function
}


func (s *Signal) Connect(f Function) {
	s.functions = append(s.functions, f)
	//fmt.Println(s.channels)
}


func (s *Signal) Emit(data interface{}) {
	for _, f := range s.functions {
		f.call(data)
	}
}
