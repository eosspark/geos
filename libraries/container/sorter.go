package container

type UInt32Slice []uint32

func (p UInt32Slice) Len() int           { return len(p) }
func (p UInt32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p UInt32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
