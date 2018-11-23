package include_test

import (
	"testing"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	. "github.com/eosspark/eos-go/plugins/appbase/app/include"
	."github.com/eosspark/eos-go/plugins/chain_interface"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
)

type Gbi struct {

}
func (g Gbi)GetBlock (s *types.SignedBlock) {
	fmt.Println("getBlock")
	fmt.Println(s.Timestamp)
}

func Test_Method(t *testing.T) {
	gbi :=App().GetMethod(GetBlockById)

	//register
	gbi.Register(&RejectedBlockFunc{Gbi{}.GetBlock})

	sb :=new(types.SignedBlock)
	sb.Timestamp = types.NewBlockTimeStamp(100)
	//CallMethods
	gbi.CallMethods(sb)


}