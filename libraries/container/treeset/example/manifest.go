package example

import "github.com/eosspark/eos-go/libraries/container"

var StringComparator = func(a, b interface{}) int { return container.StringComparator(a.(string), b.(string)) }

//go:generate go install "github.com/eosspark/eos-go/libraries/container"
//go:generate go install "github.com/eosspark/eos-go/libraries/container/redblacktree"
//go:generate go install "github.com/eosspark/eos-go/libraries/container/treeset"
//go:generate gotemplate "github.com/eosspark/eos-go/libraries/container/treeset" StringSet(string,StringComparator,false)
//go:generate gotemplate "github.com/eosspark/eos-go/libraries/container/treeset" MultiStringSet(string,StringComparator,true)
//go:generate go build .
