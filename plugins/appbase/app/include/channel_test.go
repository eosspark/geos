package include_test

import (
	"testing"
	."github.com/eosspark/eos-go/plugins/chain_interface"
	"github.com/eosspark/eos-go/chain/types"
	"fmt"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"time"
	."github.com/eosspark/eos-go/plugins/appbase/app/include"
	."github.com/eosspark/eos-go/plugins/appbase/app"
)




//var  acceptedBlockHeader *Signal;
//var  acceptedBlock *Signal;
//var  irreversibleBlock *Signal;
//var  acceptedTransaction *Signal;
//var  appliedTransaction *Signal;
//var  acceptedConfirmation *Signal;
//var  badAlloc *Signal;

type blockAcceptor struct {

}

func (*blockAcceptor) doAccept(s *types.SignedBlock) {
	fmt.Println(s.Timestamp)
}

func doAccept(s *types.SignedBlock) {
	fmt.Println(s.Timestamp)
}


func (*blockAcceptor)doRejectedBlockFunc(s *types.SignedBlock){
	fmt.Println(s.Timestamp)
}

func Test_Channel(t *testing.T) {

	//subscribe
	App().GetChannel(PreAcceptedBlock).Subscribe(&PreAcceptedBlockFunc{doAccept})
	App().GetChannel(PreAcceptedBlock).Subscribe(&PreAcceptedBlockFunc{new(blockAcceptor).doAccept})
	rbf :=&RejectedBlockFunc{Func:new(blockAcceptor).doRejectedBlockFunc}
	App().GetChannel(RejectedBlock).Subscribe(rbf)

	//call
	sb := &types.SignedBlock{}
	sb.Timestamp = types.NewBlockTimeStamp(100)
	//App().GetChannel(chain.PreAcceptedBlock).Publish(sb)
	App().GetChannel(RejectedBlock).Publish(sb)
	App().GetChannel(PreAcceptedBlock).Publish(sb)
	App().GetChannel(AcceptedBlockHeader).Publish(sb)


	timer := asio.NewDeadlineTimer(App().GetIoService())
	timer.ExpiresFromNow(time.Millisecond)
	timer.AsyncWait(func(err error) {
		App().Quit()
	})

	App().Exec()
}
