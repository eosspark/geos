package example

import "github.com/eosspark/eos-go/common/container"

var IntComparator = func(a, b interface{}) int { return container.IntComparator(a.(int), b.(int)) }
var StringComparator = func(a, b interface{}) int { return container.StringComparator(a.(string), b.(string)) }

//go:generate go install "github.com/eosspark/eos-go/common/container/treemap"
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" IntStringMap(int,string,IntComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" MultiIntStringMap(int,string,IntComparator,true)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" IntStringPtrMap(int,*string,IntComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" MultiIntStringPtrMap(int,*string,IntComparator,true)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" StringIntMap(string,int,StringComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treemap" MultiStringIntMap(string,int,StringComparator,true)
//go:generate go build .
