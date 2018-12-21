Install
-------

Install using go get

    go get github.com/ncw/gotemplate/...

and this will build the `gotemplate` binary in `$GOPATH/bin`.

And make sure install this package

    go install github.com/eosspark/container/...

Using templates
---------------

To use a template, first you must tell `gotemplate` that you want to
use it using a special comment in your code.  For example

    //go:generate gotemplate "github.com/eosspark/container/template/treeset" MySet(string,utils.StringComparator)

