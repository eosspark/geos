package example

import "github.com/eosspark/eos-go/common/container"

var StringComparator = func(a, b interface{}) int { return container.StringComparator(a.(string), b.(string)) }
//go:generate go install "github.com/eosspark/eos-go/common/container/redblacktree/"
//go:generate go install "github.com/eosspark/eos-go/common/container/treeset/"
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treeset" StringSet(string,StringComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/treeset" MultiStringSet(string,StringComparator,true)
//go:generate go build .
