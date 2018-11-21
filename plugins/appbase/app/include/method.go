package include

type Method struct {
	signal *Signal
}

func NewMethod () *Method{
	m := &Method{&Signal{}}
	return m
}
func (m *Method) Register(f Function){
	m.signal.functions= append(m.signal.functions, f)
}

func (m *Method) CallMethods(data interface{})  {
	for _,f := range m.signal.functions {
		f.call(data)
	}
}
