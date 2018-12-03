// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasmgo

import "errors"

// ErrUnreachable is the error value used while trapping the VM when
// an unreachable operator is reached during execution.
var ErrUnreachable = errors.New("wasmgo: reached unreachable")

func (vm *VM) unreachable() {
	panic(ErrUnreachable)
}

func (vm *VM) nop() {}
