// Package deps contains the console JavaScript dependencies Go embedded.
package deps

//go:generate go-bindata -nometadata -pkg deps -o bindata.go eosgo.js
//go:generate gofmt -w -s bindata.go
