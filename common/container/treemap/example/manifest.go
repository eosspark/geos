package example

//go:generate go install "github.com/eosspark/eos-go/common/container/treemap/..."
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" IntStringMap(int,string,utils.IntComparator)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" IntStringPtrMap(int,*string,utils.IntComparator)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" StringIntMap(string,int,utils.StringComparator)
//go:generate go build .
func main() {}