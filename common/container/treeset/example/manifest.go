package example

//go:generate go install "github.com/eosspark/eos-go/common/container/treeset/..."
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treeset" StringSet(string,utils.StringComparator,false)
//go:generate go build .
func main() {
}
