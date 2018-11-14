package common

//The result will be 0 if a==b, -1 if first < second, and +1 if first > second.
type ElementObject interface {
	ElementObject()
}

type Bucket struct {
	Data    []ElementObject
	Compare func(first ElementObject, second ElementObject) int
}

func (m *Bucket) Len() int {
	return len(m.Data)
}

func (m *Bucket) GetData(i int) ElementObject {
	if len(m.Data)-1 >= i {
		return m.Data[i]
	}
	return nil
}

func (m *Bucket) Clear() {
	if len(m.Data) > 0 {
		m.Data = nil
	}
}

func (m *Bucket) Find(element ElementObject) (bool, int) {
	r := m.searchSub(element)
	if r >= 0 && m.Compare(element, m.Data[r]) == 0 {
		return true, r
	}
	return false, -1
}

func (b *Bucket) Easer(element ElementObject) bool {
	result := false
	if b.Len() == 0 {
		return result
	}
	_, sub := b.Find(element)
	if sub >= 0 /* && f.Len()>=1*/ {
		b.Data = append(b.Data[:sub], b.Data[sub+1:]...)
		result = true
	}
	return result
}

func (m *Bucket) searchSub(obj ElementObject) int {
	length := m.Len()
	if length == 0 {
		return -1
	}
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			if m.Compare(m.Data[h], obj) == -1 {
				i = h + 1
			} else if m.Compare(m.Data[h], obj) == 0 {
				return h
			} else {
				j = h
			}
		}
	}
	return i
}

func (m *Bucket) Insert(obj ElementObject) (*ElementObject, bool) {
	//fmt.Println("Bucket", m.Compare == nil)
	var result ElementObject
	length := m.Len()
	target := m.Data
	exist := false
	if length == 0 {
		m.Data = append(m.Data, obj)
		result = m.Data[0]
	} else {
		r := m.searchSub(obj)
		start := m.Compare(target[0], obj)
		end := m.Compare(obj, target[length-1])
		if (start == -1 || start == 0) && (end == -1 || end == 0) {
			//Insert middle
			elemnts := []ElementObject{}
			first := target[:r]
			second := target[r:length]
			elemnts = append(elemnts, first...)
			elemnts = append(elemnts, obj)
			elemnts = append(elemnts, second...)
			m.Data = elemnts
			result = elemnts[r]
		} else {
			//insert target before
			if m.Compare(obj, target[0]) == -1 {
				elemnts := []ElementObject{}
				elemnts = append(elemnts, obj)
				elemnts = append(elemnts, target...)
				m.Data = elemnts
				result = elemnts[0]
			} else if m.Compare(obj, target[length-1]) == 1 { //target append
				target = append(target, obj)
				result = target[length]
				m.Data = target
			}
		}
	}
	return &result, exist
}
