package p2p

import (
	"fmt"
	"testing"
)

func TestSyncManger(t *testing.T) {
	sync_manager := NewsyncManager(100)
	fmt.Println(stageStr(sync_manager.state))
	sync_manager.setStage(libCatchup)
	fmt.Println(stageStr(sync_manager.state))

	fmt.Println(sync_manager.syncRequired())

}
