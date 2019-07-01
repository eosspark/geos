package include

type Method struct {
	signal *Signal
}

func NewMethod() *Method {
	m := &Method{&Signal{}}
	return m
}
func (m *Method) Register(f Caller) {
	m.signal.Connect(f)
}

func (m *Method) CallMethods(data ...interface{}) {
	m.signal.Emit(data...)
}
