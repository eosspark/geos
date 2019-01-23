package example

//go:generate go install "github.com/eosspark/eos-go/common/container/treemap"
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" IntStringMap(int,string,utils.IntComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" MultiIntStringMap(int,string,utils.IntComparator,true)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" IntStringPtrMap(int,*string,utils.IntComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" MultiIntStringPtrMap(int,*string,utils.IntComparator,true)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" StringIntMap(string,int,utils.StringComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" MultiStringIntMap(string,int,utils.StringComparator,true)
//go:generate go build .
func main() {}
