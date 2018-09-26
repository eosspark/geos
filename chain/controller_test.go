/*
 *  @Time : 2018/8/29 下午5:47
 *  @Author : xueyahui
 *  @File : controller_test.go
 *  @Software: GoLand
 */

package chain

import (
	"fmt"
	"testing"
)

func TestPopBlock(t *testing.T) {
	con := NewController()
	//con.PopBlock()
	fmt.Println(con)
}

func TestAbortBlock(t *testing.T) {
	con := NewController()

	con.AbortBlock()
	fmt.Println(con)
}

func TestSetApplayHandler(t *testing.T) {
	con := NewController()
	fmt.Println(con)
	applyCon := ApplyContext{}
	con.SetApplayHandler(111, 111, 111, applyCon)
}
